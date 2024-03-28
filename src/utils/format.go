package utils

import (
	"errors"
	"regexp"
)

// FormatCPF formata o CPF informado
func FormatCPF(cpf string) (string, error) {
	re := regexp.MustCompile("\\D")
	cpf = re.ReplaceAllString(cpf, "")

	if len(cpf) != 11 {
		return "", errors.New("CPF must have exactly 11 digits")
	}

	formattedCPF := cpf[:3] + "." + cpf[3:6] + "." + cpf[6:9] + "-" + cpf[9:]

	return formattedCPF, nil
}

// FormatTelephone formata o telefone informado
func FormatTelephone(telephone string) (string, error) {
	re := regexp.MustCompile("\\D")
	telephone = re.ReplaceAllString(telephone, "")

	if len(telephone) != 11 {
		return "", errors.New("Telephone must have exactly 11 digits")
	}
	formatTelephone := "(" + telephone[:2] + ") " + telephone[2:7] + "-" + telephone[7:]

	return formatTelephone, nil
}
