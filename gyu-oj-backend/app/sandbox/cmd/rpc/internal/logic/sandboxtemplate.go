package logic

import (
	"github.com/pkg/errors"
	"github.com/zeromicro/go-zero/core/logx"
	"gyu-oj-backend/app/sandbox/cmd/rpc/pb"
	"gyu-oj-backend/app/sandbox/models/enums"
	"gyu-oj-backend/common/xerr"
	"time"
)

var (
	GoBinaryFileName = "main"
	UserCodesDir     = "userCodes"
	TimeoutLimit     = 4000 * time.Millisecond // 时间限制（MS）
	MemoryLimit      = 128                     //内存限制（MB）
)

func SandboxTemplate(c ExecuteCodeItf, param *pb.ExecuteCodeReq) (*pb.ExecuteCodeResp, error) {
	resp := &pb.ExecuteCodeResp{}
	// 1，保存文件
	userCodePath, err := c.SaveCodeToFile([]byte(param.Code))
	if err != nil {
		logx.Infof("保存文件错误, err: %v", err)
		resp.Status = enums.SystemError.GetStatus()
		resp.Message = enums.SystemError.GetMsg()
		return resp, err
	}

	// 5，删除文件
	defer func() {
		err := c.DropFile(userCodePath)
		if err != nil {
			logx.Infof("删除文件失败: %v", err)
			return
		}
		logx.Info("删除文件成功")
	}()

	// 2，编译代码
	err = c.CompileCode(userCodePath)
	if err != nil {
		logx.Infof("编译文件错误: %v", err)
		resp.Status = enums.CompileFail.GetStatus()
		resp.Message = enums.CompileFail.GetMsg()
		return resp, err
	}

	// 3，运行代码
	runResultList, err := c.RunCode(userCodePath, param.InputList)
	if err != nil {
		logx.Infof("运行用户代码文件错误: %v", err)
		causeErr := errors.Cause(err)
		runCodeErr, _ := causeErr.(*xerr.CodeError)
		if runCodeErr.GetErrCode() == xerr.RunTimeoutError {
			resp.Status = enums.RunTimeout.GetStatus()
			resp.Message = enums.RunTimeout.GetMsg()
			return resp, err
		}
		if runCodeErr.GetErrCode() == xerr.RunOutOfMemoryError {
			resp.Status = enums.RunOutOfMemory.GetStatus()
			resp.Message = enums.RunOutOfMemory.GetMsg()
			return resp, err
		}
		resp.Status = enums.RunFail.GetStatus()
		resp.Message = enums.RunFail.GetMsg()
		return resp, err
	}
	if len(runResultList) <= 0 {
		resp.Status = enums.RunFail.GetStatus()
		resp.Message = enums.RunFail.GetMsg()
		return resp, err
	}

	logx.Info("运行代码成功")

	// 4，整理代码的输出结果
	resp = c.GetOutputResponse(runResultList)

	return resp, nil
}
