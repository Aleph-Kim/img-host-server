package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"img-host-server/internal/utils"

	"github.com/joho/godotenv"
	"golang.org/x/crypto/bcrypt"
)

// User 구조체
type User struct {
	Username string `json:"username"`
}

// 사용자 등록 (POST /users)
func SaveUser(w http.ResponseWriter, r *http.Request) {
	// .env 파일 로딩
	if err := godotenv.Load(); err != nil {
		log.Println("Error loading .env file")
		utils.RespondJSON(w, http.StatusInternalServerError, "서버 환경 설정을 불러오는 데 실패했습니다.")
		return
	}

	// 환경 변수에서 관리자 비밀번호 가져오기
	adminPassword := os.Getenv("ADMIN_PASSWORD")
	if adminPassword == "" {
		log.Println("관리자 비밀번호가 설정되지 않았습니다.")
		utils.RespondJSON(w, http.StatusInternalServerError, "관리자 비밀번호가 설정되지 않았습니다.")
		return
	}

	// 관리자 인증 체크
	authHeader := r.Header.Get("Authorization")
	if authHeader != adminPassword {
		log.Println("관리자 인증 실패")
		utils.RespondJSON(w, http.StatusForbidden, "관리자만 사용자 등록이 가능합니다.")
		return
	}

	// 요청 바디에서 username 디코딩
	var newUser User
	if err := json.NewDecoder(r.Body).Decode(&newUser); err != nil {
		log.Println("업로드 폴더 생성 실패:", err)
		utils.RespondJSON(w, http.StatusBadRequest, "잘못된 요청입니다.")
		return
	}

	// 랜덤 비밀번호(16바이트) 생성
	randomSecret, err := utils.GenerateRandomSecret(16)
	if err != nil {
		log.Println("업로드 폴더 생성 실패:", err)
		utils.RespondJSON(w, http.StatusInternalServerError, "비밀번호 생성에 실패했습니다.")
		return
	}

	// 생성된 비밀번호 bcrypt 해싱
	hashedSecret, err := bcrypt.GenerateFromPassword([]byte(randomSecret), bcrypt.DefaultCost)
	if err != nil {
		log.Println("업로드 폴더 생성 실패:", err)
		utils.RespondJSON(w, http.StatusInternalServerError, "비밀번호 해싱에 실패했습니다.")
		return
	}

	// 기존 사용자 정보 불러오기
	users, err := utils.LoadUsers()
	if err != nil {
		log.Println("업로드 폴더 생성 실패:", err)
		utils.RespondJSON(w, http.StatusInternalServerError, "사용자 정보 불러오기에 실패했습니다.")
		return
	}

	// 이미 존재하는 사용자면 에러 반환
	if _, exists := users[newUser.Username]; exists {
		log.Println("업로드 폴더 생성 실패:", err)
		utils.RespondJSON(w, http.StatusConflict, "이미 존재하는 사용자명입니다.")
		return
	}

	// 사용자 추가
	users[newUser.Username] = string(hashedSecret)

	// 업데이트된 사용자 정보 저장
	if err := utils.SaveUsers(users); err != nil {
		log.Println("업로드 폴더 생성 실패:", err)
		utils.RespondJSON(w, http.StatusInternalServerError, "파일 저장에 실패했습니다.")
		return
	}

	// 성공 응답 반환
	w.WriteHeader(http.StatusCreated)

	utils.RespondJSON(w, http.StatusCreated, fmt.Sprintf("사용자 등록 성공: Secret Key: %s", randomSecret))
}
