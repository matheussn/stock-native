package services

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"stock/internal/models"
)

var ErrProductVariationNotFound = errors.New("product variation not found")

type ProductVariationService struct {
	db *sql.DB
}

func NewProductVariationService(db *sql.DB) *ProductVariationService {
	return &ProductVariationService{db: db}
}

func (s *ProductVariationService) Create(ctx context.Context, productID int64, description string, baseQuantity int64) (models.ProductVariation, error) {
	description = strings.TrimSpace(description)

	if productID <= 0 {
		return models.ProductVariation{}, errors.New("product_id must be greater than zero")
	}
	if baseQuantity <= 0 {
		return models.ProductVariation{}, errors.New("base_quantity must be greater than zero")
	}

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return models.ProductVariation{}, fmt.Errorf("begin transaction: %w", err)
	}
	defer func() {
		_ = tx.Rollback()
	}()

	res, err := tx.ExecContext(
		ctx,
		`INSERT INTO product_variation (product_id, description, base_quantity, current_stock, is_active) VALUES (?, ?, ?, 0, 1)`,
		productID,
		nullableText(description),
		baseQuantity,
	)
	if err != nil {
		return models.ProductVariation{}, fmt.Errorf("insert product variation: %w", err)
	}

	id, err := res.LastInsertId()
	if err != nil {
		return models.ProductVariation{}, fmt.Errorf("read inserted id: %w", err)
	}

	variation, err := getProductVariationByIDTx(ctx, tx, id)
	if err != nil {
		return models.ProductVariation{}, err
	}

	if err := tx.Commit(); err != nil {
		return models.ProductVariation{}, fmt.Errorf("commit transaction: %w", err)
	}
	return variation, nil
}

func (s *ProductVariationService) ListByProduct(ctx context.Context, productID int64, includeInactive bool) ([]models.ProductVariation, error) {
	query := `
SELECT
	id,
	product_id,
	COALESCE(description, ''),
	base_quantity,
	current_stock,
	is_active,
	created_at
FROM product_variation
WHERE product_id = ?
`
	if !includeInactive {
		query += "AND is_active = 1\n"
	}
	query += "ORDER BY base_quantity DESC, id ASC"

	rows, err := s.db.QueryContext(ctx, query, productID)
	if err != nil {
		return nil, fmt.Errorf("list product variations: %w", err)
	}
	defer rows.Close()

	result := make([]models.ProductVariation, 0)
	for rows.Next() {
		var v models.ProductVariation
		if err := rows.Scan(&v.ID, &v.ProductID, &v.Description, &v.BaseQuantity, &v.CurrentStock, &v.IsActive, &v.CreatedAt); err != nil {
			return nil, fmt.Errorf("scan product variation: %w", err)
		}
		result = append(result, v)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate product variations: %w", err)
	}

	return result, nil
}

func (s *ProductVariationService) Update(ctx context.Context, id int64, description string, baseQuantity int64) (models.ProductVariation, error) {
	description = strings.TrimSpace(description)
	if baseQuantity <= 0 {
		return models.ProductVariation{}, errors.New("base_quantity must be greater than zero")
	}

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return models.ProductVariation{}, fmt.Errorf("begin transaction: %w", err)
	}
	defer func() {
		_ = tx.Rollback()
	}()

	res, err := tx.ExecContext(
		ctx,
		`UPDATE product_variation SET description = ?, base_quantity = ? WHERE id = ?`,
		nullableText(description),
		baseQuantity,
		id,
	)
	if err != nil {
		return models.ProductVariation{}, fmt.Errorf("update product variation: %w", err)
	}

	affected, err := res.RowsAffected()
	if err != nil {
		return models.ProductVariation{}, fmt.Errorf("read updated rows: %w", err)
	}
	if affected == 0 {
		return models.ProductVariation{}, ErrProductVariationNotFound
	}

	variation, err := getProductVariationByIDTx(ctx, tx, id)
	if err != nil {
		return models.ProductVariation{}, err
	}

	if err := tx.Commit(); err != nil {
		return models.ProductVariation{}, fmt.Errorf("commit transaction: %w", err)
	}
	return variation, nil
}

func (s *ProductVariationService) SetActive(ctx context.Context, id int64, isActive bool) error {
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

	res, err := tx.ExecContext(ctx, `UPDATE product_variation SET is_active = ? WHERE id = ?`, activeValue, id)
	if err != nil {
		return fmt.Errorf("update product variation active flag: %w", err)
	}

	affected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("read updated rows: %w", err)
	}
	if affected == 0 {
		return ErrProductVariationNotFound
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit transaction: %w", err)
	}
	return nil
}

func getProductVariationByIDTx(ctx context.Context, tx *sql.Tx, id int64) (models.ProductVariation, error) {
	var v models.ProductVariation
	err := tx.QueryRowContext(
		ctx,
		`SELECT id, product_id, COALESCE(description, ''), base_quantity, current_stock, is_active, created_at FROM product_variation WHERE id = ?`,
		id,
	).Scan(&v.ID, &v.ProductID, &v.Description, &v.BaseQuantity, &v.CurrentStock, &v.IsActive, &v.CreatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.ProductVariation{}, ErrProductVariationNotFound
		}
		return models.ProductVariation{}, fmt.Errorf("get product variation by id: %w", err)
	}
	return v, nil
}
