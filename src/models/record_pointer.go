package models

import (
	"errors"
	"time"
)

const (
	edição   = "edição"
	cadastro = "cadastro"
)

// RegisterPointer responsável pelo campos do registro de ponto
type RegisterPointer struct {
	ID              uint64     `json:"id,omitempty"`
	EmployeeCode    uint64     `json:"codigo_funcionario,omitempty"`
	RecordType      *string    `json:"tipo_registro,omitempty"`
	EntryTime       *time.Time `json:"hora_entrada"`
	DepartureTime   *time.Time `json:"hora_saida"`
	CriadoEm        *time.Time `json:"criadoEm,omitempty"`
	LunchEntryTime  *time.Time `json:"hora_entrada_almoco"`
	ReturnTimeLunch *time.Time `json:"hora_retorno_almoco"`
	WorkedHours     *string    `json:"horas_trabalhadas"`
	Overtime        *string    `json:"horas_extras"`
}

func (r *RegisterPointer) PreparePoint(stage string) error {
	if err := r.validate(stage); err != nil {
		return err
	}
	return nil
}

func (r *RegisterPointer) validate(stage string) error {
	if stage == edição {
		if r.RecordType == nil || *r.RecordType == "" {
			return errors.New("The record type is required and cannot be blank")
		}
	} else if r.EmployeeCode == 0 {
		return errors.New("The employee code is required and cannot be blank")
	}

	return nil
}
