package handlers

import (
	"fmt"
	"img-host-server/internal/utils"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gorilla/mux"
)

// 파일 업로드 저장 기본 경로
const uploadPath = "./uploads"

// 파일 업로드 (POST /files)
func UploadFile(w http.ResponseWriter, r *http.Request) {
	// 인증 체크
	username, err := utils.CheckAuth(r)
	if err != nil {
		utils.RespondJSON(w, http.StatusUnauthorized, err.Error())
		return
	}

	if err := r.ParseMultipartForm(10 << 20); err != nil { // 10MB 제한
		utils.RespondJSON(w, http.StatusBadRequest, "잘못된 요청입니다.")
		return
	}

	file, header, err := r.FormFile("file")
	if err != nil {
		log.Println("파일 업로드 실패:", err)
		utils.RespondJSON(w, http.StatusBadRequest, "파일을 읽어오는데 실패했습니다.")
		return
	}
	defer file.Close()

	// 유저 전용 업로드 폴더 생성 (없을 경우)
	userUploadPath := filepath.Join(uploadPath, username)
	if err := os.MkdirAll(userUploadPath, os.ModePerm); err != nil {
		log.Println("업로드 폴더 생성 실패:", err)
		utils.RespondJSON(w, http.StatusInternalServerError, "업로드 폴더 생성에 실패했습니다.")
		return
	}

	// 파일명 안전 처리
	filename := filepath.Clean(filepath.Base(header.Filename))
	// 요청 body에 filename이 있으면 해당 값을 사용, 없으면 header.Filename 사용
	if formFilename := r.FormValue("filename"); formFilename != "" {
		filename = formFilename
	}

	// 특수문자 필터링
	if !utils.IsValidFileName(filename) {
		utils.RespondJSON(w, http.StatusBadRequest, "파일명은 한글, 알파벳, 숫자, 밑줄(_), 대시(-), 점(.) 외의 문자를 허용하지 않습니다.")
		return
	}

	savePath := filepath.Join(userUploadPath, filename)

	// 같은 이름의 파일 존재 여부 체크
	if _, err := os.Stat(savePath); err == nil {
		utils.RespondJSON(w, http.StatusConflict, "같은 이름의 파일이 이미 존재합니다.")
		return
	} else if !os.IsNotExist(err) {
		log.Println("파일 존재 여부 체크 실패:", err)
		utils.RespondJSON(w, http.StatusInternalServerError, "파일 존재 여부 체크에 실패했습니다.")
		return
	}

	// 파일 저장
	if err := utils.SaveUploadedFile(file, savePath); err != nil {
		log.Println("파일 저장 실패:", err)
		utils.RespondJSON(w, http.StatusInternalServerError, "파일 저장에 실패했습니다.")
		return
	}

	utils.RespondJSON(w, http.StatusCreated, fmt.Sprintf("파일 업로드 성공: %s", filename))
}

// 파일 수정 (PUT /files/{filename})
func UpdateFile(w http.ResponseWriter, r *http.Request) {
	// 인증 체크
	username, err := utils.CheckAuth(r)
	if err != nil {
		utils.RespondJSON(w, http.StatusUnauthorized, err.Error())
		return
	}

	vars := mux.Vars(r)
	filename := vars["filename"]

	if filename == "" {
		utils.RespondJSON(w, http.StatusBadRequest, "파일 이름을 입력해주세요.")
		return
	}

	// 유저 디렉토리에서 파일 경로 지정
	filePath := filepath.Join(uploadPath, username, filepath.Clean(filename))

	// 수정할 파일이 존재하는지 확인합니다.
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		utils.RespondJSON(w, http.StatusNotFound, "수정할 파일이 존재하지 않습니다.")
		return
	}

	// 요청 파싱 (10MB 제한)
	if err := r.ParseMultipartForm(10 << 20); err != nil {
		utils.RespondJSON(w, http.StatusBadRequest, "잘못된 요청입니다.")
		return
	}

	// 새 파일 데이터를 가져옵니다.
	file, _, err := r.FormFile("file")
	if err != nil {
		log.Println("파일 업로드 실패:", err)
		utils.RespondJSON(w, http.StatusBadRequest, "파일을 읽어오는데 실패했습니다.")
		return
	}
	defer file.Close()

	// 파일 저장
	if err := utils.SaveUploadedFile(file, filePath); err != nil {
		log.Println("파일 저장 실패:", err)
		utils.RespondJSON(w, http.StatusInternalServerError, "파일 저장에 실패했습니다.")
		return
	}

	utils.RespondJSON(w, http.StatusOK, fmt.Sprintf("파일 수정 성공: %s", filename))
}

// 파일 다운로드 (GET /files/{username}/{filename})
func GetFile(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	username := vars["username"]
	filename := vars["filename"]
	size := r.URL.Query().Get("size") // size 파라미터 가져오기

	if username == "" || filename == "" {
		utils.RespondJSON(w, http.StatusBadRequest, "파일 경로를 입력해주세요.")
		return
	}

	filePath := filepath.Join(uploadPath, username, filepath.Clean(filename))

	// 파일 존재 여부 확인
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		utils.RespondJSON(w, http.StatusNotFound, "파일을 찾을 수 없습니다.")
		return
	}

	if size == "" {
		http.ServeFile(w, r, filePath)
		return
	} else {
		resizedImage, err := utils.ResizeImage(size, filePath)
		if err != nil {
			utils.RespondJSON(w, http.StatusBadRequest, err.Error())
			return
		}

		// Content-Type 설정
		w.Header().Set("Content-Type", "image/jpeg")
		w.Write(resizedImage)
	}
}

// 파일 삭제 (DELETE /files/{filename})
func DeleteFile(w http.ResponseWriter, r *http.Request) {
	// 인증 체크
	username, err := utils.CheckAuth(r)
	if err != nil {
		utils.RespondJSON(w, http.StatusUnauthorized, err.Error())
		return
	}

	vars := mux.Vars(r)
	filename := vars["filename"]

	if filename == "" {
		utils.RespondJSON(w, http.StatusBadRequest, "파일 이름을 입력해주세요.")
		return
	}

	filePath := filepath.Join(uploadPath, username, filepath.Clean(filename))

	// 파일 존재 여부 확인
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		utils.RespondJSON(w, http.StatusNotFound, "파일이 존재하지 않습니다.")
		return
	}

	// 파일 삭제
	if err := os.Remove(filePath); err != nil {
		log.Println("파일 삭제 실패:", err)
		utils.RespondJSON(w, http.StatusInternalServerError, "파일 삭제에 실패했습니다.")
		return
	}

	utils.RespondJSON(w, http.StatusOK, fmt.Sprintf("파일 삭제 성공: %s", filename))
}
