package repository

import (
	"check-point/src/models"
	"context"
	"database/sql"
	"errors"
	"log"
	"strconv"
	"strings"
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

func (e employee) ListAllEmployee(ctx context.Context) ([]models.Employee, error) {
	query := "SELECT id, nome, email, telefone, cargo, idade, cpf, criadoem FROM public.funcionario ORDER BY id ASC"

	rows, err := e.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	employees := []models.Employee{}

	for rows.Next() {
		var employee models.Employee
		if err = rows.Scan(
			&employee.ID,
			&employee.Name,
			&employee.Email,
			&employee.Telephone,
			&employee.Office,
			&employee.Age,
			&employee.CPF,
			&employee.CriadoEm); err != nil {
			return nil, err
		}
		employees = append(employees, employee)
	}
	return employees, nil

}

func (e employee) ListRepositoryParamsEmployee(ctx context.Context, params string) ([]models.Employee, error) {
	selectQuery := "SELECT id, nome, email, telefone, cargo, idade, cpf, criadoem FROM public.funcionario ORDER BY id ASC WHERE"

	condition := []string{}

	values := []interface{}{}

	searchValues := strings.Fields(params)

	for _, val := range searchValues {
		condition = append(condition, "(nome ILIKE $"+strconv.Itoa(len(values)+1)+
			" OR email ILIKE $"+strconv.Itoa(len(values)+2)+
			" OR cargo ILIKE $"+strconv.Itoa(len(values)+3)+")")
		values = append(values, "%"+val+"%", "%"+val+"%", "%"+val+"%")
	}

	if len(condition) == 0 { //len() comprimento da string, array e slice
		return nil, errors.New("No search parameters provided")
	}

	fullQuery := selectQuery + " " + strings.Join(condition, " AND ")
	rows, err := e.db.QueryContext(ctx, fullQuery, values...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	employees := []models.Employee{}

	for rows.Next() {
		var employee models.Employee
		if err := rows.Scan(
			&employee.ID,
			&employee.Name,
			&employee.Email,
			&employee.Telephone,
			&employee.Office,
			&employee.Age,
			&employee.CPF,
			&employee.CriadoEm,
		); err != nil {
			return nil, err
		}
		employees = append(employees, employee)
	}

	return employees, nil
}

func (e employee) UpdateRepositoryEmployee(ctx context.Context, ID uint64, employee models.Employee) error {
	query := "UPDATE public.funcionario SET nome = $1, email = $2, telefone = $3, cargo = $4, idade = $5 WHERE id = $6"

	_, err := e.db.ExecContext(ctx, query, employee.Name, employee.Email, employee.Telephone, employee.Office, employee.Age, ID)
	if err != nil {
		log.Printf("Error updating employee with ID %d: %v", ID, err)
		return err
	}

	return nil
}

func (e employee) DeleteRepositoryEmployee(ctx context.Context, ID uint64) error {
	query := "DELETE FROM public.funcionario WHERE id = $1"

	_, err := e.db.ExecContext(ctx, query, ID)
	if err != nil {
		log.Printf("Error delete employee with ID %d: %v", ID, err)
		return err
	}

	return nil
}

func (e employee) ListIDRepositoryEmployee(ctx context.Context, ID uint64) (models.Employee, error) {
	query := "SELECT id, nome, email, telefone, idade, cpf, cargo, criadoem FROM public.funcionario WHERE id = $1"

	var employee models.Employee
	err := e.db.QueryRowContext(ctx, query, ID).
		Scan(&employee.ID,
			&employee.Name,
			&employee.Email,
			&employee.Telephone,
			&employee.Age,
			&employee.CPF,
			&employee.Office,
			&employee.CriadoEm)
	if err != nil {
		log.Printf("Error list employee with ID %d: %v", ID, err)
		return models.Employee{}, err
	}

	return employee, nil
}
