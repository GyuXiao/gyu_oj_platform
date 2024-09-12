package logic

import (
	"context"
	"fmt"
	"gyu-oj-backend/app/sandbox/cmd/rpc/internal/svc"
	"gyu-oj-backend/app/sandbox/cmd/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type ExecuteCodeLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewExecuteCodeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ExecuteCodeLogic {
	return &ExecuteCodeLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *ExecuteCodeLogic) ExecuteCode(in *pb.ExecuteCodeReq) (*pb.ExecuteCodeResp, error) {
	// 1,new 一个 Go 原生实现的代码沙箱
	sandboxByGoNative := NewSandboxByGoNative()

	// 2,使用代码沙箱
	resp, err := SandboxTemplate(sandboxByGoNative, in)
	if err != nil {
		fmt.Println("保存编译运行代码错误：", err)
		return nil, err
	}

	//// 统计有效样例总数
	//cnt := 0
	//for _, v := range resp.OutputList {
	//	if v != "" {
	//		cnt++
	//	}
	//}
	//logc.Infof(l.ctx, "输出样例的总数: %v", cnt)
	// 正常代码
	//"package main\n\nimport \"fmt\"\n\nfunc main() {\n\tvar a, b int\n\tfmt.Scanln(&a, &b)\n\tfmt.Println(SumOfTwoNumbers(a, b))\n}\n\nfunc SumOfTwoNumbers(a, b int) int {\n\t// 解题代码请写于此处：\n\treturn a + b\n}"
	// 超时代码
	//"package main\n\nimport \"fmt\"\n\nfunc main() {\n\tvar a, b int\n\tfmt.Scanln(&a, &b)\n\tfmt.Println(SumOfTwoNumbers(a, b))\n}\n\nfunc SumOfTwoNumbers(a, b int) int {\n\t// 解题代码请写于此处：\n\t// 存在死循环\n\tfor {\n\t\ta++\n\t}\n\treturn a + b\n}\n"

	// 3,返回代码输出结果
	return resp, nil
}
