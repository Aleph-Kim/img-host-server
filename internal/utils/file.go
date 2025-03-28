package utils

import (
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/h2non/bimg"
)

// 파일 데이터를 savePath 위치에 저장
func SaveUploadedFile(file io.Reader, savePath string) error {
	out, err := os.Create(savePath)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, file)
	return err
}

/*
resizeImage는 이미지를 주어진 크기로 리사이징합니다.

Parameters:
  - size: "width x height" 형식의 문자열 (예: "100x100")
  - filePath: 리사이징할 이미지 파일 경로

Returns:
  - []byte: 리사이징된 이미지 바이트
  - error: 에러 발생 시 에러 객체 반환
*/
func ResizeImage(size string, filePath string) ([]byte, error) {
	// 이미지 파일 확장자 체크
	if !strings.HasSuffix(strings.ToLower(filePath), ".jpg") &&
		!strings.HasSuffix(strings.ToLower(filePath), ".jpeg") &&
		!strings.HasSuffix(strings.ToLower(filePath), ".png") &&
		!strings.HasSuffix(strings.ToLower(filePath), ".webp") {
		return nil, fmt.Errorf("지원하지 않는 이미지 형식입니다")
	}

	// size 파라미터 파싱
	sizes := strings.Split(size, "x")
	if len(sizes) != 2 {
		return nil, fmt.Errorf("잘못된 크기 형식입니다. (예: 100x100)")
	}

	width, err := strconv.Atoi(sizes[0])
	if err != nil {
		return nil, fmt.Errorf("잘못된 너비 값입니다")
	}

	height, err := strconv.Atoi(sizes[1])
	if err != nil {
		return nil, fmt.Errorf("잘못된 높이 값입니다")
	}

	buffer, err := bimg.Read(filePath)
	if err != nil {
		log.Println("이미지 읽기 실패:", err)
		return nil, fmt.Errorf("이미지 처리 중 오류가 발생했습니다")
	}

	newImage, err := bimg.NewImage(buffer).Resize(width, height)
	if err != nil {
		log.Println("이미지 리사이징 실패:", err)
		return nil, fmt.Errorf("이미지 처리 중 오류가 발생했습니다")
	}

	return newImage, nil
}
