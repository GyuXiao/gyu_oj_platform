package logic

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/stdcopy"
	"github.com/zeromicro/go-zero/core/logc"
	"gyu-oj-backend/app/sandbox/cmd/rpc/pb"
	"gyu-oj-backend/app/sandbox/models"
	"gyu-oj-backend/app/sandbox/models/enums"
	"gyu-oj-backend/common/tools"
	"gyu-oj-backend/common/xerr"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

/*
* 基于 docker 运行用户代码的所有步骤:
  1,把用户的代码保存为文件
  2,编译代码，得到 Go 可执行文件
  3,把编译好的文件上传到容器环境内
  4,在容器中执行代码，得到输出结果
  5,收集整理输出结果
  6,文件清理，释放空间
  7,错误处理，提升程序健壮性
*/

var (
	globalContainerID string
)

const (
	RunGoImage       = "alpine:latest"
	RunCmdStr        = "./main "
	ContainerWorkDir = "/app"
	UserCodeName     = "main.go"
)

type SandboxByDocker struct {
	Ctx context.Context
	Cli *client.Client
}

func NewSandboxByDocker(ctx context.Context, cli *client.Client) *SandboxByDocker {
	return &SandboxByDocker{
		Ctx: ctx,
		Cli: cli,
	}
}

// 1,保存用户代码文件

func (g *SandboxByDocker) SaveCodeToFile(userCode []byte) (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		logc.Infof(ctx, "获取当前文件夹目录错误: ", err)
		return "", err
	}
	// 创建存放代码文件的目录文件
	path := fmt.Sprintf("%s/%s", dir, UserCodesDir)
	_, err = os.Stat(path)
	if os.IsNotExist(err) {
		err = os.Mkdir(path, os.ModePerm)
		if err != nil {
			logc.Infof(ctx, "创建存放用户代码文件夹错误: ", err)
			return "", err
		}
	}
	// 每个用户的代码文件夹目录
	singleCodeParentPath := fmt.Sprintf("%s/%s", path, tools.GetUUID())
	err = os.Mkdir(singleCodeParentPath, os.ModePerm)
	if err != nil {
		return "", err
	}
	// 每个用户的代码的文件路径
	codePath := fmt.Sprintf("%s/%s", singleCodeParentPath, UserCodeName)
	// 创建 main.go 文件
	err = os.WriteFile(codePath, userCode, 0644)
	if err != nil {
		logc.Infof(g.Ctx, "创建用户代码文件失败: %v", err)
		return "", err
	}

	return codePath, nil
}

// 2,在 linux 本地环境里编译为可执行文件

func (g *SandboxByDocker) CompileCode(userCodePath string) error {
	parentPath := filepath.Dir(userCodePath)
	compileCmdStr := fmt.Sprintf("go build -o %s/%s %s", parentPath, GoBinaryFileName, userCodePath)
	fmt.Println("编译代码的命令: ", compileCmdStr)

	compileParts := strings.Split(compileCmdStr, " ")
	compileCmd := exec.Command(compileParts[0], compileParts[1:]...)
	var out, stderr bytes.Buffer
	compileCmd.Stderr = &stderr
	compileCmd.Stdout = &out
	// 编译成功的话，将得到可执行文件
	err := compileCmd.Run()
	if err != nil {
		logc.Infof(ctx, "编译失败: %v", err)
		return err
	}
	// 修改可执行文件的文件权限
	compileFilePath := fmt.Sprintf("%s/%s", parentPath, GoBinaryFileName)
	err = os.Chmod(compileFilePath, os.ModePerm)
	if err != nil {
		logc.Infof(ctx, "修改可执行文件的权限失败: %v", err)
		return err
	}

	logc.Infof(ctx, "编译成功: %v", out.String())
	return nil
}

// 3,运行用户代码
// 3.1,创建并启动一个容器
// 3.2,在容器中运行代码
// 3.3,安全控制（超时控制，内存监控）

func (g *SandboxByDocker) RunCode(userCodePath string, inputList []string) ([]*models.ExecResult, error) {
	// 挂载卷
	localToContainerVolume := fmt.Sprintf("%s:%s", filepath.Dir(userCodePath), ContainerWorkDir)
	// 创建并启动容器
	containerId, err := g.CreateAndStartContainer(RunGoImage, localToContainerVolume)
	if err != nil {
		return nil, err
	}
	globalContainerID = containerId

	// 结果列表
	executeResultList := make([]*models.ExecResult, len(inputList))

	doneOfWatchMemory := make(chan struct{})
	doneOfRunCode := make(chan error, 1)
	mxMemoryCh := make(chan uint64, 1)

	// 内存监控
	go func(client *client.Client, ctx context.Context, cid string) {
		mx := uint64(0)
		for {
			select {
			case <-doneOfWatchMemory:
				mxMemoryCh <- mx
				close(doneOfWatchMemory)
				return
			default:
				// 执行 docker stats 获取资源使用情况
				containerStatsResp, err := client.ContainerStats(ctx, cid, false)
				if err != nil {
					logc.Infof(ctx, "开启内存监控失败: %v", err)
					return
				}
				var statsResp container.StatsResponse
				err = json.NewDecoder(containerStatsResp.Body).Decode(&statsResp)
				if err != nil {
					logc.Infof(ctx, "json 解析失败: %v", err)
				}
				mx = max(mx, statsResp.MemoryStats.Usage)
				containerStatsResp.Body.Close()
			}
		}
	}(g.Cli, g.Ctx, containerId)

	// 运行代码
	go func(cid string) {
		// for 循环跑所有输入样例
		for i, input := range inputList {
			cmdStr := RunCmdStr + strings.TrimSpace(input)
			runCmd := strings.Split(strings.TrimSpace(cmdStr), " ")
			res, err := g.runCodeInContainer(cid, runCmd)
			if err != nil {
				logc.Infof(g.Ctx, "运行代码失败: %v", err)
				doneOfRunCode <- err
				return
			}
			// 记录答案
			executeResultList[i] = res
		}
		doneOfWatchMemory <- struct{}{}
		doneOfRunCode <- nil
		return
	}(containerId)

	select {
	case <-time.After(TimeoutLimit):
		return nil, xerr.NewErrCode(xerr.RunTimeoutError)
	case err := <-doneOfRunCode:
		if err != nil {
			return nil, err
		}
		logc.Infof(g.Ctx, "运行代码成功")
		memory, err := g.getMemoryUsage(mxMemoryCh)
		if err != nil {
			return nil, err
		}
		for i := range executeResultList {
			executeResultList[i].Memory = memory
		}
		return executeResultList, nil
	}
}

func (g *SandboxByDocker) getMemoryUsage(memoryCh chan uint64) (int64, error) {
	memory := <-memoryCh
	close(memoryCh)
	logc.Infof(ctx, "内存使用情况: %v", memory)
	if tools.BToMb(memory) > uint64(MemoryLimit) {
		return 0, xerr.NewErrCode(xerr.RunOutOfMemoryError)
	}
	return int64(tools.BToMb(memory)), nil
}

func (g *SandboxByDocker) runCodeInContainer(containerId string, runCmd []string) (*models.ExecResult, error) {
	result := &models.ExecResult{}

	// 执行 docker exec 运行代码
	startTime := time.Now()
	execCreateResp, err := g.Cli.ContainerExecCreate(g.Ctx, containerId, container.ExecOptions{
		AttachStdin:  true,
		AttachStdout: true,
		AttachStderr: true,
		WorkingDir:   ContainerWorkDir,
		Cmd:          runCmd,
		Tty:          true,
	})
	if err != nil {
		logc.Infof(g.Ctx, "执行 docker exec 失败: %v", err)
		return nil, err
	}

	// 建立连接，获取输出
	resp, err := g.Cli.ContainerExecAttach(ctx, execCreateResp.ID, container.ExecStartOptions{})
	defer resp.Close()
	if err != nil {
		logc.Infof(g.Ctx, "获取 docker 输出失败: %v", err)
		return nil, err
	}
	needTime := time.Since(startTime).Milliseconds()
	result.Time = needTime

	// 收集标准输出、标准错误
	var outBuf, errBuf bytes.Buffer
	_, err = stdcopy.StdCopy(&outBuf, &errBuf, resp.Reader)
	stdout, err := io.ReadAll(&outBuf)
	if err != nil {
		logc.Infof(g.Ctx, "读取标准输出或标准错误失败: %v", err)
		return nil, err
	}
	result.StdOut = string(stdout)

	stderr, err := io.ReadAll(&errBuf)
	if err != nil {
		logc.Infof(g.Ctx, "读取标准错误失败: %v", err)
		return nil, err
	}
	result.StdErr = string(stderr)

	return result, nil
}

// 4,整理输出数据

func (g *SandboxByDocker) GetOutputResponse(executeResult []*models.ExecResult) *pb.ExecuteCodeResp {
	resp := &pb.ExecuteCodeResp{
		Message: enums.Success.GetMsg(),
		Status:  enums.Success.GetStatus(),
	}

	outputList := make([]string, len(executeResult))
	for i, result := range executeResult {
		if result == nil {
			continue
		}
		// 如果代码运行存在错误
		if result.StdErr != "" {
			resp.Status = enums.RunFail.GetStatus()
			resp.Message = enums.RunFail.GetMsg()
			resp.ExecuteResultMessage = result.StdErr
			break
		}
		// 去掉换行符
		if strings.HasSuffix(result.StdOut, "\n") || strings.HasSuffix(result.StdOut, "\r") {
			result.StdOut = strings.TrimSuffix(result.StdOut, "\n")
			result.StdOut = strings.TrimSuffix(result.StdOut, "\r")
		}
		// 记录代码运行的输出
		outputList[i] = result.StdOut
		// 取最大消耗时间
		resp.ExecuteResultTime = max(resp.ExecuteResultTime, result.Time)
		// 取最大消耗内存
		resp.ExecuteResultMemory = max(resp.ExecuteResultMemory, result.Memory)
	}
	resp.OutputList = outputList
	// 记录最大消耗时间和最大消耗内存
	return resp
}

// 5,清理整个过程中产生的中间文件

func (g *SandboxByDocker) DropFile(userCodePath string) error {
	err := os.RemoveAll(filepath.Dir(userCodePath))
	if err != nil {
		return err
	}
	return nil
}

// 创建并启动一个 docker 容器

func (g *SandboxByDocker) CreateAndStartContainer(image, volume string) (string, error) {
	// 容器配置
	containerConfig := &container.Config{
		Image:           image,
		Tty:             true,
		AttachStdin:     true,
		AttachStdout:    true,
		AttachStderr:    true,
		WorkingDir:      ContainerWorkDir,
		NetworkDisabled: true,
	}
	hostConfig := &container.HostConfig{
		Binds: []string{volume},
		Resources: container.Resources{
			Memory: int64(MemoryLimit * 1024 * 1024),
		},
		ReadonlyRootfs: true,
	}
	// 容器名
	containerName := fmt.Sprintf("Container-%s", tools.GetUUID()[:12])
	// 创建容器
	createContainerResp, err := g.Cli.ContainerCreate(g.Ctx,
		containerConfig,
		hostConfig,
		nil,
		nil,
		containerName)
	if err != nil {
		logc.Infof(g.Ctx, "创建容器错误: %v", err)
	}

	// 容器 ID
	containerId := createContainerResp.ID
	// 启动容器
	err = g.Cli.ContainerStart(g.Ctx, containerId, container.StartOptions{})
	if err != nil {
		logc.Infof(g.Ctx, "启动容器错误: %v", err)
	}

	return containerId, nil
}

// 关停并删除一个 docker 容器

func (g *SandboxByDocker) StopAndRemoveContainer(containerId string) error {
	// 删除容器（先停后删）
	err := g.Cli.ContainerStop(ctx, containerId, container.StopOptions{})
	if err != nil {
		logc.Infof(ctx, "停止容器错误: %v", err)
		return err
	}
	err = g.Cli.ContainerRemove(ctx, containerId, container.RemoveOptions{})
	if err != nil {
		logc.Infof(ctx, "删除容器错误: %v", err)
		return err
	}
	return nil
}
