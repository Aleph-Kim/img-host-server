package utils

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
)

// 요청 헤더의 X-Username과 X-Secret 값을 이용하여 인증
func CheckAuth(r *http.Request) (string, error) {
	username := r.Header.Get("X-Username")
	secret := r.Header.Get("X-Secret")
	if username == "" || secret == "" {
		return "", errors.New("인증 정보가 부족합니다.")
	}

	// users.json 파일 열기
	file, err := os.Open("./internal/db/users.json")
	if err != nil {
		return "", errors.New("인증 정보 불러오기에 실패했습니다.")
	}
	defer file.Close()

	// users.json 파일을 파싱합니다.
	var users map[string]string
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&users); err != nil {
		return "", fmt.Errorf("users.json 파일 파싱에 실패했습니다: %v", err)
	}

	// 인증 정보 확인
	if expectedSecret, ok := users[username]; !ok || expectedSecret != secret {
		return "", errors.New("인증에 실패했습니다.")
	}
	return username, nil
}
