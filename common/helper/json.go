package helper

import (
	"encoding/json"
	"github.com/zeromicro/go-zero/core/logx"
)

type ReturnContentStruct struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Date    string `json:"date"`
}

func (r *ReturnContentStruct) JsonArrPush(jsonData string) string {
	if len(jsonData) > 0 {
		var rc ReturnContentStruct
		var rcList []ReturnContentStruct
		err := json.Unmarshal([]byte(jsonData), &rc)
		if err != nil {
			logx.Errorf("jsonArrPush:%v", err)
		}
		rcList = append(rcList, rc, *r)
		contentByte, _ := json.Marshal(rcList)
		return string(contentByte)
	} else {
		returnContentByte, _ := json.Marshal(r)
		return string(returnContentByte)
	}
}
