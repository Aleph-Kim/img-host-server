package utils

import (
	"io"
	"os"
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
