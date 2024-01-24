package models

import (
	"check-point/src/security"
	"check-point/src/utils"
	"errors"
	"strings"
	"time"

	"github.com/badoux/checkmail"
)

type Employee struct {
	ID        uint64    `json:"id,omitempty"`
	Name      string    `json:"nome,omitempty"`
	Email     string    `json:"email,omitempty"`
	Telephone string    `json:"telefone,omitempty"`
	PassWord  string    `json:"senha,omitempty"`
	CPF       string    `json:"cpf,omitempty"`
	Office    string    `json:"office,omitempty"`
	Age       uint64    `json:"age,omitempty"`
	CriadoEm  time.Time `json:"criadoEm,omitempty"`
}

func (e *Employee) Prepare(stage string) error {
	if err := e.validate(stage); err != nil {
		return err
	}
	if err := e.formatting(stage); err != nil {
		return err
	}
	return nil
}
func (e *Employee) validate(stage string) error {
	if e.Name == "" {
		return errors.New("The name is required and cannot be blank")
	}
	if e.Email == "" {
		return errors.New("Email is required and cannot be blank")
	}
	if err := checkmail.ValidateFormat(e.Email); err != nil {
		return errors.New("The email entered is invalid")
	}
	if e.Telephone == "" {
		return errors.New("Phone is required and cannot be blank")
	}
	if e.CPF == "" {
		return errors.New("CPF is required and cannot be blank")
	} else {
		err := utils.FormatCPF(e.CPF)
		if err != nil {
			return err
		}
	}
	if stage == "cadastro" && e.PassWord == "" {
		return errors.New("Password is required and cannot be blank")
	}

	return nil
}
func (u *Employee) formatting(stage string) error {
	u.Name = strings.TrimSpace(u.Name)
	u.Email = strings.TrimSpace(u.Email)
	u.Telephone = strings.TrimSpace(u.Telephone)

	if stage == "cadastro" {
		passWordHash, err := security.Hash(u.PassWord)
		if err != nil {
			return err
		}
		u.PassWord = string(passWordHash)
	}
	return nil
}
