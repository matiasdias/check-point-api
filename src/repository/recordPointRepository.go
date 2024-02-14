package repository

import (
	"check-point/src/models"
	"context"
	"database/sql"
	"errors"
	"strconv"
	"strings"
)

type recordPoint struct {
	db *sql.DB
}

// Recebe um banco que vai ser chamada lá no controlllers
func NewRepositoryRecordPoint(db *sql.DB) *recordPoint {
	return &recordPoint{db: db}
}

func (r *recordPoint) CreateRecordPoint(ctx context.Context, point models.RegisterPointer) (*models.RegisterPointer, error) {
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

func (r recordPoint) ListAllRecordPoint(ctx context.Context) ([]models.RecordWithEmployee, error) {
	query := "SELECT r.id, r.codigo_funcionario, f.nome, f.email, r.tipo_registro, f.cargo FROM public.registro_ponto r INNER JOIN public.funcionario f ON r.codigo_funcionario = f.id ORDER BY r.id ASC"

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
			&recordP.Name,
			&recordP.Email,
			&recordP.RecordType,
			&recordP.Office,
		); err != nil {
			return nil, err
		}
		record = append(record, recordP)
	}
	return record, nil

}

func (r recordPoint) ListRepositoryParamsRecordPoint(ctx context.Context, params string) ([]models.RecordWithEmployee, error) {
	selectQuery := "SELECT r.id, r.codigo_funcionario, f.nome, f.email, r.tipo_registro FROM public.registro_ponto r INNER JOIN public.funcionario f ON r.codigo_funcionario = f.id"

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
			&recordP.Name,
			&recordP.Email,
			&recordP.RecordType,
		); err != nil {
			return nil, err
		}
		record = append(record, recordP)
	}

	return record, nil
}
