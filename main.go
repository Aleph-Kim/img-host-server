package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

const uploadPath = "./uploads"

func uploadHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "POST 메서드만 허용됩니다.", http.StatusMethodNotAllowed)
		return
	}
	file, header, err := r.FormFile("image")
	if err != nil {
		http.Error(w, "파일을 읽어오는데 실패했습니다.", http.StatusBadRequest)
		return
	}
	defer file.Close()

	// 업로드 폴더가 없으면 생성합니다.
	if err := os.MkdirAll(uploadPath, os.ModePerm); err != nil {
		http.Error(w, "업로드 폴더 생성에 실패했습니다.", http.StatusInternalServerError)
		return
	}

	filename := filepath.Base(header.Filename)
	out, err := os.Create(filepath.Join(uploadPath, filename))
	if err != nil {
		http.Error(w, "파일 생성에 실패했습니다.", http.StatusInternalServerError)
		return
	}
	defer out.Close()

	_, err = io.Copy(out, file)
	if err != nil {
		http.Error(w, "파일 저장에 실패했습니다.", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "파일 업로드 성공: %s", filename)
}

func imageHandler(w http.ResponseWriter, r *http.Request) {
	imageName := r.URL.Query().Get("name")
	if imageName == "" {
		http.Error(w, "이미지 이름이 제공되지 않았습니다.", http.StatusBadRequest)
		return
	}
	imagePath := filepath.Join(uploadPath, filepath.Base(imageName))
	http.ServeFile(w, r, imagePath)
}

func main() {
	http.HandleFunc("/upload", uploadHandler)
	http.HandleFunc("/image", imageHandler)

	log.Println("서버 시작: 포트 3000")
	if err := http.ListenAndServe(":3000", nil); err != nil {
		log.Fatal(err)
	}
}
