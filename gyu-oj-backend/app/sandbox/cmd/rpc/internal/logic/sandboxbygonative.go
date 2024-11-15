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
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

type SandboxByGoNative struct {
	Ctx context.Context
}

func NewSandboxByGoNative(ctx context.Context) *SandboxByGoNative {
	return &SandboxByGoNative{
		Ctx: ctx,
	}
}

func (g *SandboxByGoNative) SaveCodeToFile(userCode []byte) (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		logc.Infof(g.Ctx, "获取当前文件夹目录错误: %v", err)
		return "", err
	}
	// 创建存放代码文件的目录文件
	path := fmt.Sprintf("%s/%s", dir, UserCodesDir)
	_, err = os.Stat(path)
	if os.IsNotExist(err) {
		err = os.Mkdir(path, os.ModePerm)
		if err != nil {
			logc.Infof(g.Ctx, "创建存放总代码文件夹错误: %v", err)
			return "", err
		}
	}
	// 每个用户的代码文件夹目录
	singleCodeParentPath := fmt.Sprintf("%s/%s", path, tools.GetUUID())
	err = os.Mkdir(singleCodeParentPath, os.ModePerm)
	if err != nil {
		logc.Infof(g.Ctx, "创建存放用户代码的文件夹错误: %v", err)
		return "", err
	}
	// 每个用户的代码的文件路径
	codePath := fmt.Sprintf("%s/main.go", singleCodeParentPath)
	// 创建 main.go 文件
	f, err := os.Create(codePath)
	if err != nil {
		logc.Infof(g.Ctx, "创建用户代码文件失败: %v", err)
		return "", err
	}
	// 写入代码
	f.Write(userCode)
	defer f.Close()

	return codePath, nil
}

func (g *SandboxByGoNative) CompileCode(userCodePath string) error {
	//if runtime.GOOS == "windows" {
	//	GoBinaryFileName = GoBinaryFileName + ".exe"
	//}

	parentPath := filepath.Dir(userCodePath)
	compileCmdStr := fmt.Sprintf("go build -o %s/%s %s", parentPath, GoBinaryFileName, userCodePath)
	logc.Infof(g.Ctx, "编译命令: %v", compileCmdStr)

	compileParts := strings.Split(compileCmdStr, " ")
	compileCmd := exec.Command(compileParts[0], compileParts[1:]...)
	var out, stderr bytes.Buffer
	compileCmd.Stderr = &stderr
	compileCmd.Stdout = &out
	// 编译成功的话，将得到可执行文件
	err := compileCmd.Run()
	if err != nil {
		logc.Infof(g.Ctx, "编译失败: %v", err)
		return xerr.NewErrCode(xerr.CompileFailError)
	}
	// 修改可执行文件的文件权限
	compileFilePath := fmt.Sprintf("%s/%s", parentPath, GoBinaryFileName)
	err = os.Chmod(compileFilePath, os.ModePerm)
	if err != nil {
		logc.Infof(g.Ctx, "修改可执行文件的权限失败: %v", err)
		return err
	}

	logc.Infof(g.Ctx, "编译成功: %v", out.String())
	return nil
}

func (g *SandboxByGoNative) RunCode(userCodePath string, inputList []string) ([]*models.ExecResult, error) {
	parentPath := filepath.Dir(userCodePath)
	runCmdStr := fmt.Sprintf("%s/%s", parentPath, GoBinaryFileName)
	//if runtime.GOOS != "windows" {
	//	runCmdStr = "./" + runCmdStr
	//}

	done := make(chan error, 1)
	executeResult := make([]*models.ExecResult, len(inputList))
	go func() {
		for i, input := range inputList {
			res, err := doRun(g.Ctx, input, runCmdStr)
			if err != nil {
				done <- err
				return
			}
			executeResult[i] = res
		}
		done <- nil
		return
	}()

	select {
	case <-time.After(TimeoutLimit):
		return nil, xerr.NewErrCode(xerr.RunTimeoutError)
	case err := <-done:
		close(done)
		if err != nil {
			return nil, err
		}
		return executeResult, nil
	}
}

func doRun(ctx context.Context, input string, runCmdStr string) (*models.ExecResult, error) {
	input = strings.TrimSpace(input)

	runCmd := &exec.Cmd{}
	runCmd = exec.Command(runCmdStr, input)
	// 标准输出、标准错误
	var out, stderr bytes.Buffer
	runCmd.Stderr = &stderr
	runCmd.Stdout = &out

	// 代码运行之前的内存
	var beforeMemory runtime.MemStats
	runtime.ReadMemStats(&beforeMemory)

	// 代码运行前的时间
	startTime := time.Now()

	// 执行代码
	err := runCmd.Run()
	if err != nil {
		logc.Infof(ctx, "运行代码失败: %v", err)
		return nil, xerr.NewErrCode(xerr.RunFailError)
	}

	// 代码运行所需时间
	needTime := time.Since(startTime).Milliseconds()
	if needTime > TimeoutLimit.Milliseconds() {
		return nil, xerr.NewErrCode(xerr.RunTimeoutError)
	}

	// 代码运行之后的内存
	var afterMemory runtime.MemStats
	runtime.ReadMemStats(&afterMemory)
	needMemory := tools.BToMb(afterMemory.Alloc) - tools.BToMb(beforeMemory.Alloc)
	if int(needMemory) > MemoryLimit {
		return nil, xerr.NewErrCode(xerr.RunOutOfMemoryError)
	}

	return &models.ExecResult{
		StdOut: out.String(),
		StdErr: stderr.String(),
		Time:   needTime,
		Memory: int64(needMemory),
	}, nil
}

func (g *SandboxByGoNative) GetOutputResponse(executeResult []*models.ExecResult) *pb.ExecuteCodeResp {
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

func (g *SandboxByGoNative) DropFile(userCodePath string) error {
	err := os.RemoveAll(filepath.Dir(userCodePath))
	if err != nil {
		return err
	}
	return nil
}
