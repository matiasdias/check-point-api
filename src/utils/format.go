package utils

import (
	"errors"
	"regexp"
)

// FormatCPF formata um CPF adicionando pontos e traços
func FormatCPF(cpf string) error {
	re := regexp.MustCompile("\\D")
	cpf = re.ReplaceAllString(cpf, "")

	if len(cpf) != 11 {
		return errors.New("CPF must have exactly 11 digits")
	}
	return nil
}
