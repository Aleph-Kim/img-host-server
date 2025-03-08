package utils

import (
	"encoding/json"
	"net/http"
)

// 응답 구조체 정의
type Response struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
}

// JSON 응답을 반환하는 유틸리티 함수
func RespondJSON(w http.ResponseWriter, status int, msg string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(Response{Code: status, Msg: msg})
}
