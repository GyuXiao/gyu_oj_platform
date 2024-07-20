package tools

import (
	"fmt"
	"github.com/GyuXiao/gyu-api-sdk/sdk/response"
	"github.com/mitchellh/mapstructure"
	"github.com/zeromicro/go-zero/core/logx"
	"reflect"
	"strconv"
)

func MapConvertStruct(source map[string]string, baseRsp *response.ErrorResponse) {
	result := make(map[string]any, len(source))
	for k, v := range source {
		if k == "Code" {
			code, _ := strconv.Atoi(v)
			result[k] = code
			continue
		}
		result[k] = v
	}
	err := mapstructure.Decode(result, &baseRsp)
	if err != nil {
		logx.Errorf("map 转换为 struct 失败 %v", err)
	}
}

func StructConvertMap(obj *response.BaseResponse) map[string]string {
	var mp map[string]any
	err := mapstructure.Decode(obj.ErrorResponse, &mp)
	if err != nil {
		logx.Errorf("struct 转换为 map 失败 %v", err)
	}
	result := make(map[string]string, len(mp))
	for k, v := range mp {
		if reflect.TypeOf(v).Kind() == reflect.Map {
			v = RemoveMapString(v)
		}
		result[k] = fmt.Sprintf("%v", v)
	}
	return result
}
