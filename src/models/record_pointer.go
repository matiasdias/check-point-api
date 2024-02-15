package models

import (
	"errors"
	"time"
)

// RegisterPointer responsavel pelo campos do registro de ponto
type RegisterPointer struct {
	ID              uint64    `json:"id,omitempty"`
	EmployeeCode    uint64    `json:"codigo_funcionario" binding:"required"`
	RecordType      string    `json:"tipo_registro" binding:"required"`
	EntryTime       time.Time `json:"hora_entrada"`
	DepartureTime   time.Time `json:"hora_saida"`
	CriadoEm        time.Time `json:"criadoEm,omitempty"`
	LunchEntryTime  time.Time `json:"hora_entrada_almoco"`
	ReturnTimeLunch time.Time `json:"hora_retorno_almoco"`
	WorkedHours     float64   `json:"horas_trabalhadas"`
}

func (r *RegisterPointer) PreparePoint() error {
	if err := r.validate(); err != nil {
		return err
	}
	return nil
}

func (r *RegisterPointer) validate() error {
	if r.EmployeeCode == 0 { // obs: o zero é considerado vazio em go
		return errors.New("The employee code is required and cannot be blank")
	}
	if r.RecordType == "" {
		return errors.New("The record type is required and cannot be blank")
	}

	return nil
}
