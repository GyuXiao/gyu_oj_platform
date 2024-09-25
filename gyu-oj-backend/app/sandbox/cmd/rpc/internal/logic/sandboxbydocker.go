package logic

import (
	"bytes"
	"context"
	"fmt"
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
	"runtime"
	"strings"
	"time"
)

//var GoBinaryFileName = "main"
//var userCodesDir = "userCodes"
//var TimeOut = 4500 * time.Millisecond // 时间限制（MS）
//var MemoryLimit = 128                 //内存限制（MB）
//
//var ctx = context.Background()

type SandboxByDocker struct {
}

func NewSandboxByDocker() *SandboxByDocker {
	return &SandboxByDocker{}
}

/*
   1,把用户的代码保存为文件
   2,编译代码，得到 Go 可执行文件
   3,把编译好的文件上传到容器环境内
   4,在容器中执行代码，得到输出结果
   5,收集整理输出结果
   6,文件清理，释放空间
   7,错误处理，提升程序健壮性
*/

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
	singleCodeParentPath := fmt.Sprintf("%s/%s", path, tools.GetUUID())
	err = os.Mkdir(singleCodeParentPath, os.ModePerm)
	if err != nil {
		return "", err
	}
	// 每个用户的代码的文件路径
	codePath := fmt.Sprintf("%s/main.go", singleCodeParentPath)
	// 创建 main.go 文件
	f, err := os.Create(codePath)
	if err != nil {
		return "", err
	}
	// 写入代码
	f.Write(userCode)
	defer f.Close()

	return codePath, nil
}

func (g *SandboxByDocker) CompileCode(userCodePath string) error {
	if runtime.GOOS == "windows" {
		GoBinaryFileName = GoBinaryFileName + ".exe"
	}

	parentPath := filepath.Dir(userCodePath)
	compileCmdStr := fmt.Sprintf("go build -o %s/%s %s", parentPath, GoBinaryFileName, userCodePath)
	logc.Infof(ctx, "编译命令: %v", compileCmdStr)

	compileParts := strings.Split(compileCmdStr, " ")
	compileCmd := exec.Command(compileParts[0], compileParts[1:]...)
	var out, stderr bytes.Buffer
	compileCmd.Stderr = &stderr
	compileCmd.Stdout = &out
	// 编译成功的话，将得到可执行文件
	err := compileCmd.Run()
	if err != nil {
		logc.Infof(ctx, "编译失败: %v", err)
		return xerr.NewErrCode(xerr.CompileFailError)
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

func (g *SandboxByDocker) RunCode(userCodePath string, inputList []string) ([]*models.ExecResult, error) {
	parentPath := filepath.Dir(userCodePath)
	runCmdStr := fmt.Sprintf("%s/%s", parentPath, GoBinaryFileName)
	if runtime.GOOS != "windows" {
		runCmdStr = "./" + runCmdStr
	}

	ctx, cancel := context.WithTimeout(context.Background(), TimeOut)
	defer cancel()
	errorChan := make(chan error, 1)
	done := make(chan struct{}, 1)

	executeResult := make([]*models.ExecResult, len(inputList))
	go func() {
		defer close(errorChan)
		defer close(done)
		err := run(inputList, runCmdStr, executeResult)
		if err != nil {
			errorChan <- err
			return
		}
		done <- struct{}{}
		return
	}()

	select {
	case <-done:
		return executeResult, nil
	case err := <-errorChan:
		return nil, err
	case <-ctx.Done():
		return nil, xerr.NewErrCode(xerr.RunTimeoutError)
	}
}

func run(inputList []string, runCmdStr string, executeResult []*models.ExecResult) error {
	runCmd := &exec.Cmd{}
	for i, input := range inputList {
		runCmd = exec.Command(runCmdStr)
		// 标准输出、标准错误
		var out, stderr bytes.Buffer
		runCmd.Stderr = &stderr
		runCmd.Stdout = &out
		stdinPipe, err := runCmd.StdinPipe()
		if err != nil {
			logc.Infof(ctx, "建立读取输入管道错误: %v", err)
		}
		// 写入输入样例
		io.WriteString(stdinPipe, input+"\n")

		// 代码运行之前的内存
		var beforeMemory runtime.MemStats
		runtime.ReadMemStats(&beforeMemory)

		// 代码运行前的时间
		startTime := time.Now()

		// 执行代码
		err = runCmd.Run()
		if err != nil {
			return xerr.NewErrCode(xerr.RunFailError)
		}

		// 代码运行所需时间
		needTime := time.Since(startTime).Milliseconds()
		if needTime > TimeOut.Milliseconds() {
			return xerr.NewErrCode(xerr.RunTimeoutError)
		}

		// 代码运行之后的内存
		var afterMemory runtime.MemStats
		runtime.ReadMemStats(&afterMemory)
		needMemory := tools.BToMb(afterMemory.Alloc) - tools.BToMb(beforeMemory.Alloc)
		if int(needMemory) > MemoryLimit {
			return xerr.NewErrCode(xerr.RunOutOfMemoryError)
		}

		executeResult[i] = &models.ExecResult{
			StdOut: out.String(),
			StdErr: stderr.String(),
			Time:   needTime,
			Memory: int64(needMemory),
		}
	}
	return nil
}

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
		if strings.HasSuffix(result.StdOut, "\n") {
			result.StdOut = strings.TrimSuffix(result.StdOut, "\n")
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
