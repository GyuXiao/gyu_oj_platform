package impl

import (
	"context"
	"encoding/json"
	"github.com/mitchellh/mapstructure"
	"github.com/zeromicro/go-zero/core/logc"
	"github.com/zeromicro/go-zero/rest/httpc"
	"gyu-oj-backend/app/judge/cmd/rpc/internal/svc"
	"gyu-oj-backend/app/judge/models/types"
	"io"
	"net/http"
)

type RemoteSandbox struct {
	Ctx *svc.ServiceContext
}

func NewRemoteSandbox(ctx *svc.ServiceContext) *RemoteSandbox {
	return &RemoteSandbox{Ctx: ctx}
}

func (sb *RemoteSandbox) ExecuteCode(req *types.ExecuteCodeReq) (*types.ExecuteCodeResp, error) {
	ctx := context.Background()
	execReq := &Request{
		InputList: req.InputList,
		Code:      req.Code,
		Language:  req.Language,
	}
	execResp, err := httpc.Do(context.Background(), http.MethodPost, sb.Ctx.Config.CodeSandbox.Url, execReq)
	if err != nil {
		logc.Infof(ctx, "调用 sandbox-rpc 服务执行代码文件错误, err: %v", err)
		return nil, err
	}

	// 读取转发请求的响应内容
	respBody, err := io.ReadAll(execResp.Body)
	if err != nil {
		logc.Infof(ctx, "读取响应代码沙箱服务响应结果错误, err: %v", err)
	}
	defer execResp.Body.Close()
	var baseResponse BaseResponse
	err = json.Unmarshal(respBody, &baseResponse)
	if err != nil {
		logc.Infof(ctx, "jsonUnmarshal err: %v", err)
	}

	var resp Response
	// 将 baseResponse.Data 的数据转换为 resp
	err = mapstructure.Decode(baseResponse.Data, &resp)
	if err != nil {
		logc.Infof(ctx, "mapstructure.Decode err: %v", err)
	}

	return &types.ExecuteCodeResp{
		OutputList: resp.OutputList,
		Message:    resp.Message,
		Status:     resp.Status,
		JudgeInfo: types.JudgeInfo{
			Message: resp.ExecuteResultMessage,
			Time:    resp.ExecuteResultTime,
			Memory:  resp.ExecuteResultMemory,
		},
	}, nil
}

type Request struct {
	InputList []string `json:"inputList"`
	Code      string   `json:"code"`
	Language  string   `json:"language"`
}

type Response struct {
	OutputList           []string `json:"outputList"`
	Message              string   `json:"message"`
	Status               int64    `json:"status"`
	ExecuteResultMessage string   `json:"executeResultMessage"`
	ExecuteResultTime    int64    `json:"executeResultTime"`
	ExecuteResultMemory  int64    `json:"executeResultMemory"`
}

type BaseResponse struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}
