package main

import (
	"bytes"
	"context"
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

var ctx = context.Background()

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
	//userCode := "package main\n\nimport \"fmt\"\n\nfunc main() {\n\tvar a, b int\n\tfmt.Scanln(&a, &b)\n\tfmt.Println(SumOfTwoNumbers(a, b))\n}\n\nfunc SumOfTwoNumbers(a, b int) int {\n\t// 解题代码请写于此处：\n\treturn a + b\n}"
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
	//time.Sleep(3 * time.Second)

	// 运行用户代码
	inputList := []string{"1 1", "2 2", "3 3", "4 4", "5 5"}
	list, err := sandboxByDocker.RunCode(userCodePath, inputList)
	if err != nil {
		logc.Infof(ctx, "运行代码失败 err: %v", err)
	}
	for _, val := range list {
		fmt.Printf("%+v\n", val)
	}
}

const (
	CompileGoImage   = "golang:1.22.0-alpine"
	CompileCmd       = "go build -o main /app/main.go"
	RunGoImage       = "alpine:latest"
	RunCmdStr        = "./main "
	ContainerWorkDir = "/app"
	UserCodeName     = "main.go"
)

var GoBinaryFileName = "main"
var userCodesDir = "userCodes"
var TimeOut = 4500 * time.Millisecond // 时间限制（MS）
var MemoryLimit = 128                 //内存限制（MB）

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

// 2,在本地环境（linux）里编译为可执行文件，然后复制到容器即可。
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

// 1，创建并启动容器
// 2，for 循环里输入样例，运行代码
// 3，删除容器

func (g *SandboxByDocker) RunCode(userCodePath string, inputList []string) ([]*ExecResult, error) {
	// 挂载卷
	localToContainerVolume := fmt.Sprintf("%s:%s", filepath.Dir(userCodePath), ContainerWorkDir)
	// 创建并启动容器
	containerId, err := g.CreateAndStartContainer(RunGoImage, localToContainerVolume)
	if err != nil {
		return nil, err
	}
	fmt.Println("创建并启动运行代码的容器成功")

	executeResult := make([]*ExecResult, len(inputList))
	for i, input := range inputList {
		cmdStr := RunCmdStr + strings.TrimSpace(input)
		runCmd := strings.Split(strings.TrimSpace(cmdStr), " ")
		fmt.Println(runCmd)

		// 创建执行命令
		execCreateResp, err := g.Cli.ContainerExecCreate(g.Ctx, containerId, container.ExecOptions{
			AttachStdin:  true,
			AttachStdout: true,
			AttachStderr: true,
			WorkingDir:   ContainerWorkDir,
			Cmd:          runCmd,
			Tty:          true,
		})
		if err != nil {
			logc.Infof(g.Ctx, "创建运行命令失败: %v", err)
			return nil, err
		}
		// 执行命令
		resp, err := g.Cli.ContainerExecAttach(ctx, execCreateResp.ID, container.ExecStartOptions{})
		if err != nil {
			logc.Infof(g.Ctx, "等待运行命令执行完成时获取输出失败: %v", err)
			return nil, err
		}

		// 标准输出、标准错误
		var outBuf, errBuf bytes.Buffer
		_, err = stdcopy.StdCopy(&outBuf, &errBuf, resp.Reader)
		stdout, err := io.ReadAll(&outBuf)
		if err != nil {
			logc.Infof(g.Ctx, "读取标准输出或标准错误失败: %v", err)
			return nil, err
		}
		fmt.Println("运行结果", string(stdout))
		stderr, err := io.ReadAll(&errBuf)
		if err != nil {
			logc.Infof(g.Ctx, "读取标准错误失败: %v", err)
			return nil, err
		}
		// 记录答案
		executeResult[i] = &ExecResult{
			StdOut: string(stdout),
			StdErr: string(stderr),
		}

		//done := make(chan error, 1)
		//go func(i int, resp types.HijackedResponse) {
		//	// 标准输出、标准错误
		//	var outBuf, errBuf bytes.Buffer
		//	_, err = stdcopy.StdCopy(&outBuf, &errBuf, resp.Reader)
		//	stdout, err := io.ReadAll(&outBuf)
		//	if err != nil {
		//		logc.Infof(g.Ctx, "读取标准输出或标准错误失败: %v", err)
		//		done <- err
		//		return
		//	}
		//	fmt.Println("运行结果", string(stdout))
		//	stderr, err := io.ReadAll(&errBuf)
		//	if err != nil {
		//		logc.Infof(g.Ctx, "读取标准错误失败: %v", err)
		//		done <- err
		//	}
		//	executeResult[i] = &ExecResult{
		//		StdOut: string(stdout),
		//		StdErr: string(stderr),
		//	}
		//	done <- nil
		//}(i, resp)
		//
		//select {
		//case err = <-done:
		//	if err != nil {
		//		return nil, err
		//	}
		//}

		resp.Close()
	}

	// 删除容器（先停后删）
	err = g.Cli.ContainerStop(g.Ctx, containerId, container.StopOptions{})
	if err != nil {
		logc.Infof(g.Ctx, "停止容器错误: %v", err)
		return nil, err
	}
	err = g.Cli.ContainerRemove(g.Ctx, containerId, container.RemoveOptions{})
	if err != nil {
		logc.Infof(g.Ctx, "删除容器错误: %v", err)
		return nil, err
	}

	return executeResult, nil
}

func (g *SandboxByDocker) CreateAndStartContainer(image, volume string) (string, error) {
	// 容器配置
	containerConfig := &container.Config{
		Image:        image,
		Tty:          true,
		AttachStdin:  true,
		AttachStdout: true,
		AttachStderr: true,
		WorkingDir:   ContainerWorkDir,
	}
	hostConfig := &container.HostConfig{
		Binds: []string{volume},
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

func GetUUID() string {
	id, _ := uuid.NewV6()
	return id.String()
}

func BToMb(b uint64) uint64 {
	return b / 1024 / 1024
}
