package models

import (
	"time"
)

// RecordWithEmployee responsável pelo retorno dos campos do registro de ponto e do funcionário
type RecordWithEmployee struct {
	ID              uint64     `json:"id,omitempty"`
	EmployeeCode    uint64     `json:"codigo_funcionario" binding:"required"`
	Name            *string    `json:"nome,omitempty" binding:"required"`
	Email           *string    `json:"email,omitempty" binding:"required"`
	Office          *string    `json:"office,omitempty" binding:"required"`
	EntryTime       *time.Time `json:"hora_entrada"`
	LunchEntryTime  *time.Time `json:"hora_entrada_almoco"`
	ReturnTimeLunch *time.Time `json:"hora_retorno_almoco"`
	DepartureTime   *time.Time `json:"hora_saida"`
	WorkedHours     *string    `json:"horas_trabalhadas"`
	Overtime        *string    `json:"horas_extras"`
}
