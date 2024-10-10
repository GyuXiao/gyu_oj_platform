package main

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/stdcopy"
	"github.com/google/uuid"
	"github.com/zeromicro/go-zero/core/logc"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

type ExecResult struct {
	StdOut string
	StdErr string
	Time   int64
	Memory int64
}

type ExecuteCodeResp struct {
	OutputList           []string
	Message              string
	Status               int64
	ExecuteResultMessage string
	ExecuteResultTime    int64
	ExecuteResultMemory  int64
}

var (
	ctx               = context.Background()
	userCodesDir      = "userCodes"             // 存放用户代码的目录
	timeout           = 4000 * time.Millisecond // 时间限制（MS）
	memoryLimit       = 128 * 1024 * 1024       //内存限制（bytes）
	globalContainerID string
)

const (
	GoBinaryFileName = "main"
	RunGoImage       = "alpine:latest"
	RunCmdStr        = "./main "
	ContainerWorkDir = "/app"
	UserCodeName     = "main.go"
)

func main() {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		logc.Infof(ctx, "创建 docker client 失败 err: %v", err)
		return
	}

	defer cli.Close()

	// new 一个基于 docker 实现的代码沙箱
	sandboxByDocker := NewSandboxByDocker(ctx, cli)

	// 用户代码
	userCode := "package main\n\nimport (\n\t\"fmt\"\n\t\"os\"\n\t\"strconv\"\n)\n\nfunc main() {\n\t// 已处理输入参数，示例代码直接使用\n\targs := os.Args[1:]\n\ta, _ := strconv.Atoi(args[0])\n\tb, _ := strconv.Atoi(args[1])\n\tfmt.Println(sumOfTwoNumbers(a, b))\n}\n\nfunc sumOfTwoNumbers(a, b int) int {\n\t// 请在此处编辑代码\n\treturn a + b\n}"
	// 保存用户代码
	userCodePath, err := sandboxByDocker.SaveCodeToFile([]byte(userCode))
	if err != nil {
		logc.Infof(ctx, "保存文件失败 err: %v", err)
		return
	}
	// 编译用户代码
	err = sandboxByDocker.CompileCode(userCodePath)
	if err != nil {
		logc.Infof(ctx, "编译代码失败 err: %v", err)
	}

	// 运行用户代码
	inputList := []string{"1 1", "2 2", "3 3", "4 4", "5 5"}
	list, err := sandboxByDocker.RunCode(userCodePath, inputList)
	if err != nil {
		logc.Infof(ctx, "运行代码失败 err: %v", err)
	}

	// 整理输出结果
	resp := sandboxByDocker.GetOutputResponse(list)
	cnt := 0
	if resp != nil && len(resp.OutputList) > 0 {
		for _, v := range resp.OutputList {
			if v != "" {
				cnt++
			}
		}
	}
	logc.Infof(context.Background(), "输出样例的总数: %v", cnt)
	logc.Infof(context.Background(), "代码沙箱返回结果: %+v", resp)

	// 删除用户文件
	err = sandboxByDocker.DropFile(userCodePath)
	if err != nil {
		logc.Infof(ctx, "删除文件失败 err: %v", err)
	}
	// 删除容器
	if globalContainerID == "" {
		return
	}
	err = sandboxByDocker.StopAndRemoveContainer(globalContainerID)
	if err != nil {
		logc.Infof(ctx, "删除容器失败 err: %v", err)
	}
}

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

func (g *SandboxByDocker) SaveCodeToFile(userCode []byte) (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		logc.Infof(ctx, "获取当前文件夹目录错误: ", err)
		return "", err
	}
	// 创建存放代码文件的目录文件
	path := fmt.Sprintf("%s/%s", dir, userCodesDir)
	_, err = os.Stat(path)
	if os.IsNotExist(err) {
		err = os.Mkdir(path, os.ModePerm)
		if err != nil {
			logc.Infof(ctx, "创建存放用户代码文件夹错误: ", err)
			return "", err
		}
	}
	// 每个用户的代码文件夹目录
	singleCodeParentPath := fmt.Sprintf("%s/%s", path, GetUUID())
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

// 2,在 linux 环境里编译为可执行文件，然后复制到容器即可。
// 不需要去创建容器，删除容器

func (g *SandboxByDocker) CompileCode(userCodePath string) error {
	parentPath := filepath.Dir(userCodePath)
	compileCmdStr := fmt.Sprintf("go build -o %s/%s %s", parentPath, GoBinaryFileName, userCodePath)
	fmt.Println("编译命令: ", compileCmdStr)

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

// 1，创建并启动一个容器
// 2，for 循环里输入样例，运行代码

func (g *SandboxByDocker) RunCode(userCodePath string, inputList []string) ([]*ExecResult, error) {
	// 挂载卷
	localToContainerVolume := fmt.Sprintf("%s:%s", filepath.Dir(userCodePath), ContainerWorkDir)
	// 创建并启动容器
	containerId, err := g.CreateAndStartContainer(RunGoImage, localToContainerVolume)
	if err != nil {
		return nil, err
	}
	globalContainerID = containerId

	executeResultList := []*ExecResult{}

	doneOfMemory := make(chan struct{})
	mxMemoryCh := make(chan uint64, 1)
	defer func() {
		memory := <-mxMemoryCh
		for i := range executeResultList {
			executeResultList[i].Memory = int64(BToMb(memory))
		}
	}()

	// 内存监控
	go func(client *client.Client, ctx context.Context, cid string, done chan struct{}) {
		mx := uint64(0)
		for {
			select {
			case <-done:
				mxMemoryCh <- mx
				close(done)
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
	}(g.Cli, g.Ctx, containerId, doneOfMemory)

	// 运行代码
	for _, input := range inputList {
		cmdStr := RunCmdStr + strings.TrimSpace(input)
		runCmd := strings.Split(strings.TrimSpace(cmdStr), " ")
		res, err := g.runCodeInContainer(containerId, runCmd)
		if err != nil {
			logc.Infof(g.Ctx, "运行代码失败: %v", err)
			return nil, err
		}
		// 记录答案
		executeResultList = append(executeResultList, res)
	}
	doneOfMemory <- struct{}{}

	logc.Infof(g.Ctx, "运行代码成功")
	return executeResultList, nil
}

func (g *SandboxByDocker) runCodeInContainer(containerId string, runCmd []string) (*ExecResult, error) {
	result := &ExecResult{}

	// 运行代码
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

	resp, err := g.Cli.ContainerExecAttach(ctx, execCreateResp.ID, container.ExecStartOptions{})
	defer resp.Close()
	if err != nil {
		logc.Infof(g.Ctx, "获取 docker 输出失败: %v", err)
		return nil, err
	}
	needTime := time.Since(startTime).Milliseconds()
	result.Time = needTime

	done := make(chan error, 1)
	// 收集标准输出、标准错误
	go func(reader *bufio.Reader, result *ExecResult) {
		var outBuf, errBuf bytes.Buffer
		_, err = stdcopy.StdCopy(&outBuf, &errBuf, reader)
		stdout, err := io.ReadAll(&outBuf)
		if err != nil {
			logc.Infof(g.Ctx, "读取标准输出或标准错误失败: %v", err)
			done <- err
			return
		}
		result.StdOut = string(stdout)
		stderr, err := io.ReadAll(&errBuf)
		if err != nil {
			logc.Infof(g.Ctx, "读取标准错误失败: %v", err)
			done <- err
			return
		}
		result.StdErr = string(stderr)
		done <- nil
		return
	}(resp.Reader, result)

	select {
	case <-time.After(timeout):
		return nil, errors.New("运行代码超时")
	case err := <-done:
		if err != nil {
			return nil, err
		}
		return result, nil
	}
}

func (g *SandboxByDocker) GetOutputResponse(executeResult []*ExecResult) *ExecuteCodeResp {
	resp := &ExecuteCodeResp{
		Message: Success.GetMsg(),
		Status:  Success.GetStatus(),
	}

	outputList := make([]string, len(executeResult))
	for i, result := range executeResult {
		if result == nil {
			continue
		}
		// 如果代码运行存在错误
		if result.StdErr != "" {
			resp.Status = RunFail.GetStatus()
			resp.Message = RunFail.GetMsg()
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

func (g *SandboxByDocker) DropFile(userCodePath string) error {
	err := os.RemoveAll(filepath.Dir(userCodePath))
	if err != nil {
		logc.Infof(ctx, "删除文件错误: %v", err)
		return err
	}
	return nil
}

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
			Memory: int64(memoryLimit),
		},
		ReadonlyRootfs: true,
	}
	// 容器名
	containerName := fmt.Sprintf("Container-%s", GetUUID()[:12])
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

func GetUUID() string {
	id, _ := uuid.NewV6()
	return id.String()
}

func BToMb(b uint64) uint64 {
	return b / 1024 / 1024
}

const (
	SystemError    ExecuteStatus = -1
	Success        ExecuteStatus = 0
	CompileFail    ExecuteStatus = 1
	RunFail        ExecuteStatus = 2
	RunTimeout     ExecuteStatus = 3
	RunOutOfMemory ExecuteStatus = 4
)

var ExecuteStatusMsg = map[ExecuteStatus]string{
	Success:        "代码运行正常",
	CompileFail:    "代码编译失败",
	RunFail:        "代码运行失败",
	RunTimeout:     "代码运行超时",
	RunOutOfMemory: "代码运行所需内存超过限制",
	SystemError:    "系统错误",
}

type ExecuteStatus int64

func (es ExecuteStatus) GetStatus() int64 {
	return int64(es)
}

func (es ExecuteStatus) GetMsg() string {
	return ExecuteStatusMsg[es]
}
