package services

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"stock/internal/models"
)

var ErrFamilyGroupAssignmentNotFound = errors.New("family group assignment not found")

type FamilyGroupAssignmentService struct {
	db *sql.DB
}

func NewFamilyGroupAssignmentService(db *sql.DB) *FamilyGroupAssignmentService {
	return &FamilyGroupAssignmentService{db: db}
}

func (s *FamilyGroupAssignmentService) Assign(ctx context.Context, familyID, productGroupID int64) (models.FamilyGroupAssignment, error) {
	if familyID <= 0 {
		return models.FamilyGroupAssignment{}, errors.New("family_id must be greater than zero")
	}
	if productGroupID <= 0 {
		return models.FamilyGroupAssignment{}, errors.New("product_group_id must be greater than zero")
	}

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return models.FamilyGroupAssignment{}, fmt.Errorf("begin transaction: %w", err)
	}
	defer func() {
		_ = tx.Rollback()
	}()

	if _, err := tx.ExecContext(
		ctx,
		`UPDATE family_group_assignment SET ended_at = datetime('now') WHERE family_id = ? AND ended_at IS NULL`,
		familyID,
	); err != nil {
		return models.FamilyGroupAssignment{}, fmt.Errorf("close active assignments: %w", err)
	}

	res, err := tx.ExecContext(
		ctx,
		`INSERT INTO family_group_assignment (family_id, product_group_id) VALUES (?, ?)`,
		familyID,
		productGroupID,
	)
	if err != nil {
		return models.FamilyGroupAssignment{}, fmt.Errorf("insert family group assignment: %w", err)
	}

	id, err := res.LastInsertId()
	if err != nil {
		return models.FamilyGroupAssignment{}, fmt.Errorf("read inserted id: %w", err)
	}

	assignment, err := getFamilyGroupAssignmentByIDTx(ctx, tx, id)
	if err != nil {
		return models.FamilyGroupAssignment{}, err
	}

	if err := tx.Commit(); err != nil {
		return models.FamilyGroupAssignment{}, fmt.Errorf("commit transaction: %w", err)
	}

	return assignment, nil
}

func (s *FamilyGroupAssignmentService) ListByFamily(ctx context.Context, familyID int64) ([]models.FamilyGroupAssignment, error) {
	rows, err := s.db.QueryContext(
		ctx,
		`
SELECT
	id,
	family_id,
	product_group_id,
	started_at,
	COALESCE(ended_at, '')
FROM family_group_assignment
WHERE family_id = ?
ORDER BY started_at DESC, id DESC
`,
		familyID,
	)
	if err != nil {
		return nil, fmt.Errorf("list family group assignments: %w", err)
	}
	defer rows.Close()

	result := make([]models.FamilyGroupAssignment, 0)
	for rows.Next() {
		var assignment models.FamilyGroupAssignment
		if err := rows.Scan(
			&assignment.ID,
			&assignment.FamilyID,
			&assignment.ProductGroupID,
			&assignment.StartedAt,
			&assignment.EndedAt,
		); err != nil {
			return nil, fmt.Errorf("scan family group assignment: %w", err)
		}
		result = append(result, assignment)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate family group assignments: %w", err)
	}

	return result, nil
}

func getFamilyGroupAssignmentByIDTx(ctx context.Context, tx *sql.Tx, id int64) (models.FamilyGroupAssignment, error) {
	var assignment models.FamilyGroupAssignment
	err := tx.QueryRowContext(
		ctx,
		`SELECT id, family_id, product_group_id, started_at, COALESCE(ended_at, '') FROM family_group_assignment WHERE id = ?`,
		id,
	).Scan(
		&assignment.ID,
		&assignment.FamilyID,
		&assignment.ProductGroupID,
		&assignment.StartedAt,
		&assignment.EndedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.FamilyGroupAssignment{}, ErrFamilyGroupAssignmentNotFound
		}
		return models.FamilyGroupAssignment{}, fmt.Errorf("get family group assignment by id: %w", err)
	}
	return assignment, nil
}
