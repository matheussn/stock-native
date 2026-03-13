package services

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"stock/internal/models"
)

var ErrAssistentialWorkNotFound = errors.New("assistential work not found")

type AssistentialWorkService struct {
	db *sql.DB
}

func NewAssistentialWorkService(db *sql.DB) *AssistentialWorkService {
	return &AssistentialWorkService{db: db}
}

func (s *AssistentialWorkService) Create(ctx context.Context, name, description string) (models.AssistentialWork, error) {
	name = strings.TrimSpace(name)
	if name == "" {
		return models.AssistentialWork{}, errors.New("name is required")
	}

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return models.AssistentialWork{}, fmt.Errorf("begin transaction: %w", err)
	}
	defer func() {
		_ = tx.Rollback()
	}()

	res, err := tx.ExecContext(
		ctx,
		`INSERT INTO assistential_work (name, description, is_active) VALUES (?, ?, 1)`,
		name,
		nullableText(description),
	)
	if err != nil {
		return models.AssistentialWork{}, fmt.Errorf("insert assistential work: %w", err)
	}

	id, err := res.LastInsertId()
	if err != nil {
		return models.AssistentialWork{}, fmt.Errorf("read inserted id: %w", err)
	}

	aw, err := getAssistentialWorkByIDTx(ctx, tx, id)
	if err != nil {
		return models.AssistentialWork{}, err
	}

	if err := tx.Commit(); err != nil {
		return models.AssistentialWork{}, fmt.Errorf("commit transaction: %w", err)
	}

	return aw, nil
}

func (s *AssistentialWorkService) List(ctx context.Context, includeInactive bool) ([]models.AssistentialWork, error) {
	query := `
SELECT
	id,
	name,
	COALESCE(description, ''),
	is_active,
	created_at
FROM assistential_work
`
	if !includeInactive {
		query += "WHERE is_active = 1\n"
	}
	query += "ORDER BY id ASC"

	rows, err := s.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("list assistential works: %w", err)
	}
	defer rows.Close()

	result := make([]models.AssistentialWork, 0)
	for rows.Next() {
		var aw models.AssistentialWork
		if err := rows.Scan(&aw.ID, &aw.Name, &aw.Description, &aw.IsActive, &aw.CreatedAt); err != nil {
			return nil, fmt.Errorf("scan assistential work: %w", err)
		}
		result = append(result, aw)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate assistential works: %w", err)
	}

	return result, nil
}

func (s *AssistentialWorkService) Update(ctx context.Context, id int64, name, description string) (models.AssistentialWork, error) {
	name = strings.TrimSpace(name)
	if name == "" {
		return models.AssistentialWork{}, errors.New("name is required")
	}

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return models.AssistentialWork{}, fmt.Errorf("begin transaction: %w", err)
	}
	defer func() {
		_ = tx.Rollback()
	}()

	res, err := tx.ExecContext(
		ctx,
		`UPDATE assistential_work SET name = ?, description = ? WHERE id = ?`,
		name,
		nullableText(description),
		id,
	)
	if err != nil {
		return models.AssistentialWork{}, fmt.Errorf("update assistential work: %w", err)
	}

	affected, err := res.RowsAffected()
	if err != nil {
		return models.AssistentialWork{}, fmt.Errorf("read updated rows: %w", err)
	}
	if affected == 0 {
		return models.AssistentialWork{}, ErrAssistentialWorkNotFound
	}

	aw, err := getAssistentialWorkByIDTx(ctx, tx, id)
	if err != nil {
		return models.AssistentialWork{}, err
	}

	if err := tx.Commit(); err != nil {
		return models.AssistentialWork{}, fmt.Errorf("commit transaction: %w", err)
	}

	return aw, nil
}

func (s *AssistentialWorkService) SetActive(ctx context.Context, id int64, isActive bool) error {
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
		`UPDATE assistential_work SET is_active = ? WHERE id = ?`,
		activeValue,
		id,
	)
	if err != nil {
		return fmt.Errorf("update assistential work active flag: %w", err)
	}

	affected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("read updated rows: %w", err)
	}
	if affected == 0 {
		return ErrAssistentialWorkNotFound
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit transaction: %w", err)
	}

	return nil
}

func getAssistentialWorkByIDTx(ctx context.Context, tx *sql.Tx, id int64) (models.AssistentialWork, error) {
	var aw models.AssistentialWork
	err := tx.QueryRowContext(
		ctx,
		`SELECT id, name, COALESCE(description, ''), is_active, created_at FROM assistential_work WHERE id = ?`,
		id,
	).Scan(&aw.ID, &aw.Name, &aw.Description, &aw.IsActive, &aw.CreatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.AssistentialWork{}, ErrAssistentialWorkNotFound
		}
		return models.AssistentialWork{}, fmt.Errorf("get assistential work by id: %w", err)
	}
	return aw, nil
}

func nullableText(value string) any {
	value = strings.TrimSpace(value)
	if value == "" {
		return nil
	}
	return value
}
