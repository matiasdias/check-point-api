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

type Employee struct {
	db *sql.DB
}

// NewRepositoryEmployee Recebe um banco que vai ser chamada através do controllers
func NewRepositoryEmployee(db *sql.DB) *Employee {
	return &Employee{db: db}
}

// CreateRepositoryEmployee responsável pelo cadastro de um novo funcionário
func (e Employee) CreateRepositoryEmployee(ctx context.Context, employee models.Employee) (*models.EmployeeResponse, error) {

	insertQuery := "INSERT INTO public.funcionario (nome, email, telefone, senha, idade, cpf, cargo, is_admin) VALUES ($1, $2, $3, $4, $5, $6, $7, $8) RETURNING id"

	_, err := e.db.ExecContext(ctx, insertQuery, employee.Name, employee.Email, employee.Telephone, employee.PassWord, employee.Age, employee.CPF, employee.Office, employee.Is_Admin)
	if err != nil {
		return nil, err
	}

	var employeeID uint64
	queryID := "SELECT id FROM public.funcionario WHERE email = $1"
	err = e.db.QueryRowContext(ctx, queryID, employee.Email).Scan(&employeeID)
	if err != nil {
		return nil, err
	}

	employeeResponse := &models.EmployeeResponse{
		ID:        employeeID,
		Name:      employee.Name,
		Email:     employee.Email,
		Telephone: employee.Telephone,
		CPF:       employee.CPF,
		Office:    employee.Office,
		Age:       employee.Age,
		Is_Admin:  employee.Is_Admin,
	}

	return employeeResponse, nil
}

// ListAllEmployee responsável pela listagem dos funcionários
func (e Employee) ListAllEmployee(ctx context.Context, isAdmin bool) ([]models.Employee, error) {
	var query string
	if isAdmin {
		query = "SELECT id, nome, email, telefone, cargo, idade, cpf, is_admin, criadoem FROM public.funcionario ORDER BY id ASC"
	} else {
		query = "SELECT id, nome, email, telefone, cargo, idade, concat(substring(cpf, 1, length(cpf)-8), '********') as cpf, is_admin, criadoem FROM public.funcionario ORDER BY id ASC"
	}

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
			&employee.Is_Admin,
			&employee.CriadoEm); err != nil {
			return nil, err
		}
		employees = append(employees, employee)
	}
	return employees, nil

}

// ListRepositoryParamsEmployee responsável pela listagem do funcionário por parametros
func (e Employee) ListRepositoryParamsEmployee(ctx context.Context, params string, isAdmin bool) ([]models.Employee, error) {
	var selectQuery string
	if isAdmin {
		selectQuery = "SELECT id, nome, email, telefone, cargo, idade, cpf, is_admin, criadoem FROM public.funcionario"
	} else {
		selectQuery = "SELECT id, nome, email, telefone, cargo, idade, concat(substring(cpf, 1, length(cpf)-8), '********') as cpf, is_admin, criadoem FROM public.funcionario"
	}

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

	fullQuery := selectQuery + " WHERE " + strings.Join(condition, " AND ") + " ORDER BY id ASC"
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
			&employee.Is_Admin,
			&employee.CriadoEm,
		); err != nil {
			return nil, err
		}
		employees = append(employees, employee)
	}

	return employees, nil
}

// UpdateRepositoryEmployee responsável pela atualização do funcionário
func (e Employee) UpdateRepositoryEmployee(ctx context.Context, ID uint64, employee models.Employee) error {
	query := "UPDATE public.funcionario SET nome = $1, email = $2, telefone = $3, cargo = $4, idade = $5, updateem = current_timestamp WHERE id = $6"

	_, err := e.db.ExecContext(ctx, query,
		employee.Name,
		employee.Email,
		employee.Telephone,
		employee.Office,
		employee.Age, ID)
	if err != nil {
		log.Printf("Error updating employee with ID %d: %v", ID, err)
		return err
	}

	return nil
}

// DeleteRepositoryEmployee responsável pela exclusão do funcionǽrio
func (e Employee) DeleteRepositoryEmployee(ctx context.Context, ID uint64) error {
	query := "DELETE FROM public.funcionario WHERE id = $1"

	_, err := e.db.ExecContext(ctx, query, ID)
	if err != nil {
		log.Printf("Error delete employee with ID %d: %v", ID, err)
		return err
	}

	return nil
}

// ListIDRepositoryEmployee responsávio pela listagem do funcionǽrio por Id
func (e Employee) ListIDRepositoryEmployee(ctx context.Context, ID uint64) (models.Employee, error) {
	// não retornar o cpf por que é um dado sensivel e não deve ser exibido
	query := "SELECT id, nome, email, telefone, idade, cpf, cargo, is_admin, criadoem FROM public.funcionario WHERE id = $1"

	var employee models.Employee
	err := e.db.QueryRowContext(ctx, query, ID).
		Scan(&employee.ID,
			&employee.Name,
			&employee.Email,
			&employee.Telephone,
			&employee.Age,
			&employee.CPF,
			&employee.Office,
			&employee.Is_Admin,
			&employee.CriadoEm)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Printf("Error list employee with ID %d: %v", ID, err)
			return models.Employee{}, nil
		}
	}

	return employee, nil
}

// UpdatePassWord responsável por atualizar a senha do funcionário
func (e Employee) UpdatePassWord(ctx context.Context, employeeID uint64, passWord string) error {
	query := "UPDATE public.funcionario SET senha = $1 WHERE id = $2"

	_, err := e.db.ExecContext(ctx, query, passWord, employeeID)
	if err != nil {
		return err
	}

	return nil
}

// FindByEmail responsável por retornar o email e a senha do funcionário
func (e Employee) FindByEmail(ctx context.Context, Email string) (models.Employee, error) {
	query := "SELECT id, email, is_admin, senha FROM public.funcionario WHERE email = $1"

	var employee models.Employee
	err := e.db.QueryRowContext(ctx, query, Email).Scan(&employee.ID, &employee.Email, &employee.Is_Admin, &employee.PassWord)
	if err != nil {
		return models.Employee{}, err
	}

	return employee, nil

}

// FindByEmailExists responsável por verificar se o email existe
func (e Employee) FindByEmailExists(ctx context.Context, Email string) (bool, error) {
	query := "SELECT EXISTS (SELECT 1 FROM public.funcionario WHERE email = $1)"

	var exists bool
	err := e.db.QueryRowContext(ctx, query, Email).Scan(&exists)
	if err != nil {
		log.Printf("Error checking if email exists: %v", err)
		return false, err
	}

	return exists, nil
}
