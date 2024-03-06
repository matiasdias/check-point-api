package models

// RecordWithEmployee responsável pelo retorno dos campos do registro de ponto e do funcionário
type RecordWithEmployee struct {
	ID           uint64 `json:"id,omitempty"`
	EmployeeCode uint64 `json:"codigo_funcionario" binding:"required"`
	Name         string `json:"nome,omitempty" binding:"required"`
	Email        string `json:"email,omitempty" binding:"required"`
	Office       string `json:"office,omitempty" binding:"required"`
}
