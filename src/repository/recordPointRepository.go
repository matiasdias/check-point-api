package repository

import (
	"check-point/src/models"
	"context"
	"database/sql"
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
