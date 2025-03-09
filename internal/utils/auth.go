package utils

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"

	"crypto/rand"
	"encoding/base64"

	"golang.org/x/crypto/bcrypt"
)

const usersFilePath = "./internal/db/users.json"

// 랜덤 값을 생성, URL-safe base64로 인코딩한 문자열을 반환
func GenerateRandomSecret(n int) (string, error) {
	bytes := make([]byte, n)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(bytes), nil
}

// 파일에서 사용자 정보를 읽어와서 map으로 반환
func LoadUsers() (map[string]string, error) {
	users := make(map[string]string)
	data, err := os.ReadFile(usersFilePath)
	if err != nil {
		if os.IsNotExist(err) {
			return users, nil
		}
		return nil, err
	}
	if len(data) > 0 {
		if err := json.Unmarshal(data, &users); err != nil {
			// 파싱 실패하면 빈 map 반환
			users = make(map[string]string)
		}
	}
	return users, nil
}

// 사용자 정보 저장
func SaveUsers(users map[string]string) error {
	file, err := os.Create(usersFilePath)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	return encoder.Encode(users)
}

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

	// users.json 파일을 파싱
	var users map[string]string
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&users); err != nil {
		return "", fmt.Errorf("users.json 파일 파싱에 실패했습니다: %v", err)
	}

	// 저장된 해시값 가져오기
	hashedSecret, ok := users[username]
	if !ok {
		return "", errors.New("인증에 실패했습니다.")
	}

	// bcrypt를 사용하여 해시 비교
	if err := bcrypt.CompareHashAndPassword([]byte(hashedSecret), []byte(secret)); err != nil {
		return "", errors.New("인증에 실패했습니다.")
	}

	return username, nil
}
