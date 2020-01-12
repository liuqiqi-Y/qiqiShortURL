package main

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"gopkg.in/validator.v2"
)

// ReqForShortURL 短链接请求结构体
type ReqForShortURL struct {
	URL                     string `json:"url" validate:"nonzero"`
	ExpirationTimeInMinutes int    `json:"expiration_time_in_minutes" validate:"min=0"`
}

// RespForShortURL 短链接响应结构体
type RespForShortURL struct {
	ShortURL string `json:"short_URL"`
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
	s, err := Shorten(req.URL, req.ExpirationTimeInMinutes)
	if err != nil {
		ResponseMsg(w, BadRequest, err.Error(), "")
	} else {
		ResponseMsg(w, Success, "", RespForShortURL{ShortURL: s})
	}
}

// ShortURLInfo 获取短URL的相关信息
func ShortURLInfo(w http.ResponseWriter, r *http.Request) {
	param := r.FormValue("shortURL")
	s, err := URLInfo(param)
	if s == "" && err == nil {
		ResponseMsg(w, BadRequest, "have no this short url", "")
		return
	}
	if err != nil {
		ResponseMsg(w, BadRequest, err.Error(), "")
		return
	}
	ResponseMsg(w, Success, "", s)
}

// Redirect 重定向
func Redirect(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	s, err := ShortURLToURL(vars["shortURL"])
	if s == "" && err == nil {
		ResponseMsg(w, BadRequest, "have no this short url", "")
		return
	}
	if err != nil {
		ResponseMsg(w, BadRequest, err.Error(), "")
		return
	}
	http.Redirect(w, r, s, http.StatusTemporaryRedirect)
}
