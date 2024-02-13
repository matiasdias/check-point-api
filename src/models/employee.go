package models

import (
	"check-point/src/security"
	"check-point/src/utils"
	"errors"
	"strings"
	"time"

	"github.com/badoux/checkmail"
)

const (
	Cadastro = "cadastro"
	Edição   = "edicao"
)

type Employee struct {
	ID        uint64    `json:"id,omitempty"`
	Name      string    `json:"nome,omitempty" binding:"required"`
	Email     string    `json:"email,omitempty" binding:"required"`
	Telephone string    `json:"telefone,omitempty" binding:"required"`
	PassWord  string    `json:"senha,omitempty" binding:"required"`
	CPF       string    `json:"cpf,omitempty" binding:"required"`
	Office    string    `json:"office,omitempty" binding:"required"`
	Age       uint64    `json:"age,omitempty" binding:"required"`
	CriadoEm  time.Time `json:"criadoEm,omitempty"`
}

type EmployeeResponse struct {
	ID        uint64    `json:"id,omitempty"`
	Name      string    `json:"nome,omitempty" binding:"required"`
	Email     string    `json:"email,omitempty" binding:"required"`
	Telephone string    `json:"telefone,omitempty" binding:"required"`
	CPF       string    `json:"cpf,omitempty" binding:"required"`
	Office    string    `json:"office,omitempty" binding:"required"`
	Age       uint64    `json:"age,omitempty" binding:"required"`
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

	if stage == Cadastro {
		if e.CPF == "" {
			return errors.New("CPF is required and cannot be blank")
		} else {
			err := utils.FormatCPF(e.CPF)
			if err != nil {
				return err
			}
		}
	}
	if e.Age < 18 {
		return errors.New("Age cannot be less than 18")
	}
	if e.Office == "" {
		return errors.New("Office is required and cannot be blank")
	}
	if stage == Cadastro && e.PassWord == "" {
		return errors.New("Password is required and cannot be blank")
	}

	return nil
}

func (e *Employee) formatting(stage string) error {
	e.Name = strings.TrimSpace(e.Name)
	e.Email = strings.TrimSpace(e.Email)
	e.Telephone = strings.TrimSpace(e.Telephone)

	if stage == Cadastro {
		passWordHash, err := security.Hash(e.PassWord)
		if err != nil {
			return err
		}
		e.PassWord = string(passWordHash)
	}
	return nil
}
