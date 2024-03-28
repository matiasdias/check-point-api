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

type RecordPoint struct {
	db *sql.DB
}

// NewRepositoryRecordPoint Recebe um banco que vai ser chamada através do controllers
func NewRepositoryRecordPoint(db *sql.DB) *RecordPoint {
	return &RecordPoint{db: db}
}

func (r *RecordPoint) CreateRecordPoint(ctx context.Context, point models.RegisterPointer) (*models.RegisterPointer, error) {
	insertQuery := "INSERT INTO public.registro_ponto (codigo_funcionario, tipo_registro) VALUES ($1, $2) RETURNING id"

	_, err := r.db.ExecContext(ctx, insertQuery, point.EmployeeCode, point.RecordType)
	if err != nil {
		return nil, err
	}

	var pointID uint64
	queryID := "SELECT id FROM public.registro_ponto WHERE codigo_funcionario = $1"
	err = r.db.QueryRowContext(ctx, queryID, point.EmployeeCode).Scan(&pointID)
	if err != nil {
		return nil, err
	}

	point.ID = pointID
	return &point, nil
}

func (r RecordPoint) ListAllRecordPoint(ctx context.Context) ([]models.RecordWithEmployee, error) {
	query := `SELECT r.id, r.codigo_funcionario, r.hora_entrada, r.hora_entrada_almoco, r.hora_retorno_almoco, r.hora_saida, r.horas_trabalhadas, r.horas_extras, f.nome, f.email, f.cargo FROM public.registro_ponto r INNER JOIN public.funcionario f ON r.codigo_funcionario = f.id ORDER BY r.id ASC`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	record := []models.RecordWithEmployee{}
	for rows.Next() {
		var recordP models.RecordWithEmployee
		if err = rows.Scan(
			&recordP.ID,
			&recordP.EmployeeCode,
			&recordP.EntryTime,
			&recordP.LunchEntryTime,
			&recordP.ReturnTimeLunch,
			&recordP.DepartureTime,
			&recordP.WorkedHours,
			&recordP.Overtime,
			&recordP.Name,
			&recordP.Email,
			&recordP.Office,
		); err != nil {
			return nil, err
		}
		record = append(record, recordP)
	}
	return record, nil

}

func (r RecordPoint) ListRepositoryParamsRecordPoint(ctx context.Context, params string) ([]models.RecordWithEmployee, error) {
	selectQuery := `SELECT r.id, r.codigo_funcionario, r.hora_entrada, r.hora_entrada_almoco, 
	r.hora_retorno_almoco, r.hora_saida, r.horas_trabalhadas, r.horas_extras,  f.nome, f.email, f.cargo FROM public.registro_ponto r 
	INNER JOIN public.funcionario f ON r.codigo_funcionario = f.id`

	condition := []string{}

	values := []interface{}{}

	searchValues := strings.Fields(params)

	for _, val := range searchValues {
		condition = append(condition, "(tipo_registro ILIKE $"+strconv.Itoa(len(values)+1)+")")
		values = append(values, "%"+val+"%")
	}

	if len(condition) == 0 { //len() comprimento da string, array e slice
		return nil, errors.New("No search parameters provided")
	}

	fullQuery := selectQuery + " WHERE " + strings.Join(condition, " AND ") + " ORDER BY id ASC"
	rows, err := r.db.QueryContext(ctx, fullQuery, values...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	record := []models.RecordWithEmployee{}

	for rows.Next() {
		var recordP models.RecordWithEmployee
		if err := rows.Scan(
			&recordP.ID,
			&recordP.EmployeeCode,
			&recordP.EntryTime,
			&recordP.LunchEntryTime,
			&recordP.ReturnTimeLunch,
			&recordP.DepartureTime,
			&recordP.WorkedHours,
			&recordP.Overtime,
			&recordP.Name,
			&recordP.Email,
			&recordP.Office,
		); err != nil {
			return nil, err
		}
		record = append(record, recordP)
	}

	return record, nil
}

func (r RecordPoint) ListEmployeeIDRepositoryRecordPoint(ctx context.Context, employeeID uint64) (models.RecordWithEmployee, error) {
	query := `SELECT r.id, r.codigo_funcionario, r.hora_entrada, r.hora_entrada_almoco, 
	r.hora_retorno_almoco, r.hora_saida, r.horas_trabalhadas, r.horas_extras, f.nome, f.email, f.cargo FROM public.registro_ponto r 
	INNER JOIN public.funcionario f ON r.codigo_funcionario = f.id WHERE f.id = $1`
	var record models.RecordWithEmployee
	err := r.db.QueryRowContext(ctx, query, employeeID).Scan(
		&record.ID,
		&record.EmployeeCode,
		&record.EntryTime,
		&record.LunchEntryTime,
		&record.ReturnTimeLunch,
		&record.DepartureTime,
		&record.WorkedHours,
		&record.Overtime,
		&record.Name,
		&record.Email,
		&record.Office,
	)
	if err == sql.ErrNoRows {
		log.Printf("Error when searching time record by employee ID %d: %v", employeeID, err)
		return models.RecordWithEmployee{}, nil
	}
	return record, nil
}

func (r RecordPoint) DeleteRecordPoint(ctx context.Context, employeeID uint64) error {
	//query := "DELETE FROM public.registro_ponto r INNER JOIN public.funcionario f ON r.codigo_funcionario = f.id WHERE f.id = $1"
	query := "DELETE FROM public.registro_ponto USING public.funcionario WHERE public.registro_ponto.codigo_funcionario = public.funcionario.id AND public.funcionario.id = $1"
	_, err := r.db.ExecContext(ctx, query, employeeID)
	if err != nil {
		log.Printf("Error deleted record type with employee ID %d: %v", employeeID, err)
		return err
	}
	return nil
}

func (r RecordPoint) UpdateRecordPoint(ctx context.Context, employeeID uint64, point models.RegisterPointer) error {
	query := "UPDATE public.registro_ponto SET tipo_registro = $1 WHERE codigo_funcionario = $2"
	_, err := r.db.ExecContext(ctx, query, point.RecordType, employeeID)
	if err != nil {
		log.Printf("Error updating record type with employee ID %d: %v", employeeID, err)
		return err
	}
	return nil
}

func (r RecordPoint) FindByEmployeeExists(ctx context.Context, employeeCode uint64) (bool, error) {
	query := "SELECT EXISTS (SELECT 1 FROM public.funcionario WHERE id = $1)"

	var exists bool
	err := r.db.QueryRowContext(ctx, query, employeeCode).Scan(&exists)
	if err != nil {
		log.Printf("Error checking if id employee exists: %v", err)
		return false, err
	}

	return exists, nil
}
