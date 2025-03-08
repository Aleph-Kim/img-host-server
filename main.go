package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

const uploadPath = "./uploads"

// 응답 구조체 정의
type Response struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
}

// JSON 응답을 반환하는 헬퍼 함수
func respondJSON(w http.ResponseWriter, status int, msg string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(Response{Code: status, Msg: msg})
}

// 파일 업로드 핸들러
func uploadHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		respondJSON(w, http.StatusMethodNotAllowed, "POST 메서드만 허용됩니다.")
		return
	}

	file, header, err := r.FormFile("image")
	if err != nil {
		log.Println("파일 업로드 실패:", err)
		respondJSON(w, http.StatusBadRequest, "파일을 읽어오는데 실패했습니다.")
		return
	}
	defer file.Close()

	// 업로드 폴더 생성 (없을 경우)
	if err := os.MkdirAll(uploadPath, os.ModePerm); err != nil {
		log.Println("업로드 폴더 생성 실패:", err)
		respondJSON(w, http.StatusInternalServerError, "업로드 폴더 생성에 실패했습니다.")
		return
	}

	// 파일명 안전 처리 및 저장 경로 설정
	filename := filepath.Base(header.Filename)
	filename = filepath.Clean(filename) // 파일 경로 공격 방지
	savePath := filepath.Join(uploadPath, filename)

	out, err := os.Create(savePath)
	if err != nil {
		log.Println("파일 생성 실패:", err)
		respondJSON(w, http.StatusInternalServerError, "파일 생성에 실패했습니다.")
		return
	}
	defer out.Close()

	// 파일 복사 및 저장
	if _, err = io.Copy(out, file); err != nil {
		log.Println("파일 저장 실패:", err)
		respondJSON(w, http.StatusInternalServerError, "파일 저장에 실패했습니다.")
		return
	}

	respondJSON(w, http.StatusOK, fmt.Sprintf("파일 업로드 성공: %s", filename))
}

// 이미지 파일 제공 핸들러
func imageHandler(w http.ResponseWriter, r *http.Request) {
	imageName := r.URL.Query().Get("name")
	if imageName == "" {
		respondJSON(w, http.StatusBadRequest, "이미지 이름이 제공되지 않았습니다.")
		return
	}

	// 안전한 파일 경로 설정
	imagePath := filepath.Join(uploadPath, filepath.Clean(imageName))
	http.ServeFile(w, r, imagePath)
}

func main() {
	http.HandleFunc("/upload", uploadHandler)
	http.HandleFunc("/image", imageHandler)

	log.Println("서버 시작: 포트 3000")
	if err := http.ListenAndServe(":3000", nil); err != nil {
		log.Fatal("서버 실행 실패:", err)
	}
}
