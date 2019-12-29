package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"gopkg.in/validator.v2"
)

// 错误类型
const (
	BadRequest  = 40001
	ServerError = 50001
)

// ReqForShortURL 短链接请求结构体
type ReqForShortURL struct {
	URL                     string `json:"url" validate:"nonzero"`
	ExpirationTimeInMinutes int    `json:"expiration_time_in_minutes" validate:"min=0"`
}

// RespForShortURL 短链接响应结构体
type RespForShortURL struct {
	ShortLink string `json:"short_link"`
}

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

// InitializeRouters 初始化路由
func InitializeRouters(router *mux.Router) {
	router.HandleFunc("/api/shortURL", CreateShortURL).Methods("POST")
	router.HandleFunc("/api/shortURLInfo", ShortURLInfo).Methods("GET")
	router.HandleFunc("/{shortURL:[a-zA-Z0-9]{1,11}}", Redirect).Methods("GET")
}

// CreateShortURL 生产短URL
func CreateShortURL(w http.ResponseWriter, r *http.Request) {
	req := ReqForShortURL{}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		ResponseMsg(w, BadRequest, err.Error(), "")
		Trace.Println("failed to decode: ", err)
		return
	}
	if err := validator.Validate(req); err != nil {
		ResponseMsg(w, BadRequest, err.Error(), "")
		Trace.Println("invalid param: ", err)
		return
	}
	fmt.Printf("%v\n", req)
}

// ShortURLInfo 获取短URL的相关信息
func ShortURLInfo(w http.ResponseWriter, r *http.Request) {
	param := r.FormValue("shortURL")
	fmt.Printf("%s\n", param)
}

// Redirect 重定向
func Redirect(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	fmt.Printf("%s\n", vars["shortURL"])
}

// Run run
func Run(address string, router *mux.Router) {
	server := &http.Server{
		Addr:    address,
		Handler: router,
	}
	server.ListenAndServe()
}
