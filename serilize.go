package main

import (
	"encoding/json"
	"net/http"
)

// 错误类型
const (
	Success     = 20000
	BadRequest  = 40001
	ServerError = 50001
)

// RespMsg 用于返回响应结构体
type RespMsg struct {
	Code int
	Msg  string
	Data interface{}
}

// NewResp 返回一个响应结构体
func NewResp(code int, msg string, data interface{}) *RespMsg {
	return &RespMsg{
		Code: code,
		Msg:  msg,
		Data: data,
	}
}

// Serialize 序列化
func Serialize(data interface{}) ([]byte, error) {
	s, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}
	return s, nil
}

// ResponseMsg 返回响应信息
func ResponseMsg(w http.ResponseWriter, code int, msg string, data interface{}) {
	respMsg := NewResp(code, msg, data)
	serialization, err := Serialize(respMsg)
	if err != nil {
		w.Write([]byte("some thing wrong with server, please try later!"))
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(serialization)
}
