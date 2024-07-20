package tools

import (
	"fmt"
	"strings"
)

func RemoveMapString(s any) string {
	str := fmt.Sprintf("%v", s)
	if strings.HasPrefix(str, "map[") && strings.HasSuffix(str, "]") {
		return str[4 : len(str)-1]
	}
	return str
}
