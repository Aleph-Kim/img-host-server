// utils/sanitize.go
package utils

import "regexp"

// 파일명에서 한글, 알파벳, 숫자, 밑줄(_), 대시(-), 점(.) 외의 문자가 포함됐는지 체크
func IsValidFileName(name string) bool {
	re := regexp.MustCompile(`^[a-zA-Z0-9_.\-가-힣]+$`)
	return re.MatchString(name)
}
