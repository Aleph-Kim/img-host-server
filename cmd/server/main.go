package main

import (
	"log"
	"net/http"

	"img-host-server/internal/handlers"

	"github.com/gorilla/mux"
)

func main() {
	r := mux.NewRouter()

	// 파일 관련 RESTful 엔드포인트 설정
	r.HandleFunc("/files", handlers.UploadFile).Methods("POST")
	r.HandleFunc("/files/{filename}", handlers.UpdateFile).Methods("PUT")
	r.HandleFunc("/files/{username}/{filename}", handlers.GetFile).Methods("GET")
	r.HandleFunc("/files/{filename}", handlers.DeleteFile).Methods("DELETE")

	log.Println("서버 시작: 포트 3000")
	if err := http.ListenAndServe(":3000", r); err != nil {
		log.Fatal("서버 실행 실패:", err)
	}
}
