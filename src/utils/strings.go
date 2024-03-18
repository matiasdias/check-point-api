package utils

import (
	"regexp"
	"strings"
)

// OnlyNumbers retorna somente os número em uma string
func OnlyNumbers(valor string) string {
	re := regexp.MustCompile(`\d`)
	return strings.Join(re.FindAllString(valor, -1), "")
}
