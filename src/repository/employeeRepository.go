package repository

import (
	"check-point/src/models"
	"context"
	"database/sql"
)

type employee struct {
	db *sql.DB
}

// Recebe um banco que vai ser chamada lá no controlllers
func NewRepositoryEmployee(db *sql.DB) *employee {
	return &employee{db: db}
}

func (e employee) CreateRepositoryEmployee(ctx context.Context, employee models.Employee) (*models.Employee, error) {

	insertQuery := "INSERT INTO public.funcionario (nome, email, telefone, senha, idade, cpf, cargo) VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id"

	_, err := e.db.ExecContext(ctx, insertQuery, employee.Name, employee.Email, employee.Telephone, employee.PassWord, employee.Age, employee.CPF, employee.Office)
	if err != nil {
		return nil, err
	}

	var employeeID uint64
	queryID := "SELECT id FROM public.funcionario WHERE email = $1"
	err = e.db.QueryRowContext(ctx, queryID, employee.Email).Scan(&employeeID)
	if err != nil {
		return nil, err
	}

	employee.ID = employeeID
	return &employee, nil
}
