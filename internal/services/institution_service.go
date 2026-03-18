package services

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"stock/internal/models"
)

var ErrInstitutionNotFound = errors.New("institution not found")
var ErrInstitutionInactive = errors.New("institution is inactive")

type InstitutionService struct {
	db *sql.DB
}

func NewInstitutionService(db *sql.DB) *InstitutionService {
	return &InstitutionService{db: db}
}

func (s *InstitutionService) Create(ctx context.Context, name, address, cnpj, responsibleName, phone string) (models.Institution, error) {
	name = strings.TrimSpace(name)
	if name == "" {
		return models.Institution{}, errors.New("name is required")
	}

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return models.Institution{}, fmt.Errorf("begin transaction: %w", err)
	}
	defer func() {
		_ = tx.Rollback()
	}()

	res, err := tx.ExecContext(
		ctx,
		`INSERT INTO institution (name, address, cnpj, responsible_name, phone, is_active) VALUES (?, ?, ?, ?, ?, 1)`,
		name,
		nullableText(address),
		nullableText(cnpj),
		nullableText(responsibleName),
		nullableText(phone),
	)
	if err != nil {
		return models.Institution{}, fmt.Errorf("insert institution: %w", err)
	}

	id, err := res.LastInsertId()
	if err != nil {
		return models.Institution{}, fmt.Errorf("read inserted id: %w", err)
	}

	institution, err := getInstitutionByIDTx(ctx, tx, id)
	if err != nil {
		return models.Institution{}, err
	}

	if err := tx.Commit(); err != nil {
		return models.Institution{}, fmt.Errorf("commit transaction: %w", err)
	}

	return institution, nil
}

func (s *InstitutionService) List(ctx context.Context, includeInactive bool) ([]models.Institution, error) {
	query := `
SELECT
	id,
	name,
	COALESCE(address, ''),
	COALESCE(cnpj, ''),
	COALESCE(responsible_name, ''),
	COALESCE(phone, ''),
	is_active,
	created_at
FROM institution
`
	if !includeInactive {
		query += "WHERE is_active = 1\n"
	}
	query += "ORDER BY id ASC"

	rows, err := s.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("list institutions: %w", err)
	}
	defer rows.Close()

	result := make([]models.Institution, 0)
	for rows.Next() {
		var institution models.Institution
		if err := rows.Scan(
			&institution.ID,
			&institution.Name,
			&institution.Address,
			&institution.CNPJ,
			&institution.ResponsibleName,
			&institution.Phone,
			&institution.IsActive,
			&institution.CreatedAt,
		); err != nil {
			return nil, fmt.Errorf("scan institution: %w", err)
		}
		result = append(result, institution)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate institutions: %w", err)
	}

	return result, nil
}

func (s *InstitutionService) Update(ctx context.Context, id int64, name, address, cnpj, responsibleName, phone string) (models.Institution, error) {
	name = strings.TrimSpace(name)
	if name == "" {
		return models.Institution{}, errors.New("name is required")
	}

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return models.Institution{}, fmt.Errorf("begin transaction: %w", err)
	}
	defer func() {
		_ = tx.Rollback()
	}()

	res, err := tx.ExecContext(
		ctx,
		`UPDATE institution SET name = ?, address = ?, cnpj = ?, responsible_name = ?, phone = ? WHERE id = ?`,
		name,
		nullableText(address),
		nullableText(cnpj),
		nullableText(responsibleName),
		nullableText(phone),
		id,
	)
	if err != nil {
		return models.Institution{}, fmt.Errorf("update institution: %w", err)
	}

	affected, err := res.RowsAffected()
	if err != nil {
		return models.Institution{}, fmt.Errorf("read updated rows: %w", err)
	}
	if affected == 0 {
		return models.Institution{}, ErrInstitutionNotFound
	}

	institution, err := getInstitutionByIDTx(ctx, tx, id)
	if err != nil {
		return models.Institution{}, err
	}

	if err := tx.Commit(); err != nil {
		return models.Institution{}, fmt.Errorf("commit transaction: %w", err)
	}

	return institution, nil
}

func (s *InstitutionService) SetActive(ctx context.Context, id int64, isActive bool) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("begin transaction: %w", err)
	}
	defer func() {
		_ = tx.Rollback()
	}()

	activeValue := 0
	if isActive {
		activeValue = 1
	}

	res, err := tx.ExecContext(
		ctx,
		`UPDATE institution SET is_active = ? WHERE id = ?`,
		activeValue,
		id,
	)
	if err != nil {
		return fmt.Errorf("update institution active flag: %w", err)
	}

	affected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("read updated rows: %w", err)
	}
	if affected == 0 {
		return ErrInstitutionNotFound
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit transaction: %w", err)
	}

	return nil
}

func getInstitutionByIDTx(ctx context.Context, tx *sql.Tx, id int64) (models.Institution, error) {
	var institution models.Institution
	err := tx.QueryRowContext(
		ctx,
		`SELECT id, name, COALESCE(address, ''), COALESCE(cnpj, ''), COALESCE(responsible_name, ''), COALESCE(phone, ''), is_active, created_at FROM institution WHERE id = ?`,
		id,
	).Scan(
		&institution.ID,
		&institution.Name,
		&institution.Address,
		&institution.CNPJ,
		&institution.ResponsibleName,
		&institution.Phone,
		&institution.IsActive,
		&institution.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.Institution{}, ErrInstitutionNotFound
		}
		return models.Institution{}, fmt.Errorf("get institution by id: %w", err)
	}
	return institution, nil
}

func ensureActiveInstitutionTx(ctx context.Context, tx *sql.Tx, id int64) error {
	var isActive int64
	err := tx.QueryRowContext(ctx, `SELECT is_active FROM institution WHERE id = ?`, id).Scan(&isActive)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ErrInstitutionNotFound
		}
		return fmt.Errorf("get institution active flag: %w", err)
	}

	if isActive == 0 {
		return ErrInstitutionInactive
	}

	return nil
}
