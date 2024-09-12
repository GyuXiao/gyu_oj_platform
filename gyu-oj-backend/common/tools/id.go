package tools

import "github.com/google/uuid"

// 生成唯一码

func GetUUID() string {
	id, _ := uuid.NewV6()
	return id.String()
}
