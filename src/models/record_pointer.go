package models

import "time"

type RegisterPointer struct {
	ID              uint64    `json:"id,omitempty"`
	CodeEmployee    uint64    `json:"code_funcionario,omitempty" binding:"required"`
	RecordType      string    `json:"tipo_registro,omitempty" binding:"required"`
	EntryTime       time.Time `json:"hora_entrada"`
	DepartureTime   time.Time `json:"hora_saida"`
	CriadoEm        time.Time `json:"criadoEm,omitempty"`
	LunchEntryTime  time.Time `json:"hora_entrada_almoco"`
	ReturnTimeLunch time.Time `json:"hora_retorno_almoco"`
	WorkedHours     float64   `json:"horas_trabalhadas"`
}
