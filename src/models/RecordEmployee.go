package models

// RecordWithEmployee junção das duas structs(employee, recordPoint)
type RecordWithEmployee struct {
	ID           uint64 `json:"id,omitempty"`
	EmployeeCode uint64 `json:"codigo_funcionario" binding:"required"`
	RecordType   string `json:"tipo_registro" binding:"required"`
	Name         string `json:"nome,omitempty" binding:"required"`
	Email        string `json:"email,omitempty" binding:"required"`
	Office       string `json:"office,omitempty" binding:"required"`
}
