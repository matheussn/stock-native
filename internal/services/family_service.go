package services

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"stock/internal/models"
)

var ErrFamilyNotFound = errors.New("family not found")

type FamilyService struct {
	db *sql.DB
}

func NewFamilyService(db *sql.DB) *FamilyService {
	return &FamilyService{db: db}
}

func (s *FamilyService) Create(ctx context.Context, assistentialWorkID int64, name string, memberCount int64, address, contact string) (models.Family, error) {
	name = strings.TrimSpace(name)
	if assistentialWorkID <= 0 {
		return models.Family{}, errors.New("assistential_work_id must be greater than zero")
	}
	if name == "" {
		return models.Family{}, errors.New("name is required")
	}
	if memberCount <= 0 {
		return models.Family{}, errors.New("member_count must be greater than zero")
	}

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return models.Family{}, fmt.Errorf("begin transaction: %w", err)
	}
	defer func() {
		_ = tx.Rollback()
	}()

	res, err := tx.ExecContext(
		ctx,
		`INSERT INTO family (assistential_work_id, name, member_count, address, contact, is_active) VALUES (?, ?, ?, ?, ?, 1)`,
		assistentialWorkID,
		name,
		memberCount,
		nullableText(address),
		nullableText(contact),
	)
	if err != nil {
		return models.Family{}, fmt.Errorf("insert family: %w", err)
	}

	id, err := res.LastInsertId()
	if err != nil {
		return models.Family{}, fmt.Errorf("read inserted id: %w", err)
	}

	family, err := getFamilyByIDTx(ctx, tx, id)
	if err != nil {
		return models.Family{}, err
	}

	if err := tx.Commit(); err != nil {
		return models.Family{}, fmt.Errorf("commit transaction: %w", err)
	}

	return family, nil
}

func (s *FamilyService) List(ctx context.Context, includeInactive bool, assistentialWorkID int64) ([]models.Family, error) {
	query := `
SELECT
	id,
	assistential_work_id,
	name,
	member_count,
	COALESCE(address, ''),
	COALESCE(contact, ''),
	is_active,
	created_at
FROM family
WHERE 1 = 1
`
	args := make([]any, 0)

	if !includeInactive {
		query += "AND is_active = 1\n"
	}
	if assistentialWorkID > 0 {
		query += "AND assistential_work_id = ?\n"
		args = append(args, assistentialWorkID)
	}
	query += "ORDER BY id ASC"

	rows, err := s.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("list families: %w", err)
	}
	defer rows.Close()

	result := make([]models.Family, 0)
	for rows.Next() {
		var family models.Family
		if err := rows.Scan(
			&family.ID,
			&family.AssistentialWorkID,
			&family.Name,
			&family.MemberCount,
			&family.Address,
			&family.Contact,
			&family.IsActive,
			&family.CreatedAt,
		); err != nil {
			return nil, fmt.Errorf("scan family: %w", err)
		}
		result = append(result, family)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate families: %w", err)
	}

	return result, nil
}

func (s *FamilyService) Update(ctx context.Context, id int64, name string, memberCount int64, address, contact string) (models.Family, error) {
	name = strings.TrimSpace(name)
	if name == "" {
		return models.Family{}, errors.New("name is required")
	}
	if memberCount <= 0 {
		return models.Family{}, errors.New("member_count must be greater than zero")
	}

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return models.Family{}, fmt.Errorf("begin transaction: %w", err)
	}
	defer func() {
		_ = tx.Rollback()
	}()

	res, err := tx.ExecContext(
		ctx,
		`UPDATE family SET name = ?, member_count = ?, address = ?, contact = ? WHERE id = ?`,
		name,
		memberCount,
		nullableText(address),
		nullableText(contact),
		id,
	)
	if err != nil {
		return models.Family{}, fmt.Errorf("update family: %w", err)
	}

	affected, err := res.RowsAffected()
	if err != nil {
		return models.Family{}, fmt.Errorf("read updated rows: %w", err)
	}
	if affected == 0 {
		return models.Family{}, ErrFamilyNotFound
	}

	family, err := getFamilyByIDTx(ctx, tx, id)
	if err != nil {
		return models.Family{}, err
	}

	if err := tx.Commit(); err != nil {
		return models.Family{}, fmt.Errorf("commit transaction: %w", err)
	}
	return family, nil
}

func (s *FamilyService) SetActive(ctx context.Context, id int64, isActive bool) error {
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

	res, err := tx.ExecContext(ctx, `UPDATE family SET is_active = ? WHERE id = ?`, activeValue, id)
	if err != nil {
		return fmt.Errorf("update family active flag: %w", err)
	}

	affected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("read updated rows: %w", err)
	}
	if affected == 0 {
		return ErrFamilyNotFound
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit transaction: %w", err)
	}

	return nil
}

func getFamilyByIDTx(ctx context.Context, tx *sql.Tx, id int64) (models.Family, error) {
	var family models.Family
	err := tx.QueryRowContext(
		ctx,
		`SELECT id, assistential_work_id, name, member_count, COALESCE(address, ''), COALESCE(contact, ''), is_active, created_at FROM family WHERE id = ?`,
		id,
	).Scan(
		&family.ID,
		&family.AssistentialWorkID,
		&family.Name,
		&family.MemberCount,
		&family.Address,
		&family.Contact,
		&family.IsActive,
		&family.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.Family{}, ErrFamilyNotFound
		}
		return models.Family{}, fmt.Errorf("get family by id: %w", err)
	}
	return family, nil
}
