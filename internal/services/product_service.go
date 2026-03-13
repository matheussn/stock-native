package services

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"stock/internal/models"
)

var ErrProductNotFound = errors.New("product not found")

type ProductService struct {
	db *sql.DB
}

func NewProductService(db *sql.DB) *ProductService {
	return &ProductService{db: db}
}

func (s *ProductService) Create(ctx context.Context, name, baseUnit, description string) (models.Product, error) {
	name = strings.TrimSpace(name)
	baseUnit = strings.TrimSpace(baseUnit)

	if name == "" {
		return models.Product{}, errors.New("name is required")
	}
	if !isValidBaseUnit(baseUnit) {
		return models.Product{}, errors.New("base_unit must be one of: g, kg, ml, L, un")
	}

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return models.Product{}, fmt.Errorf("begin transaction: %w", err)
	}
	defer func() {
		_ = tx.Rollback()
	}()

	res, err := tx.ExecContext(
		ctx,
		`INSERT INTO product (name, base_unit, description, is_active) VALUES (?, ?, ?, 1)`,
		name,
		baseUnit,
		nullableText(description),
	)
	if err != nil {
		return models.Product{}, fmt.Errorf("insert product: %w", err)
	}

	id, err := res.LastInsertId()
	if err != nil {
		return models.Product{}, fmt.Errorf("read inserted id: %w", err)
	}

	product, err := getProductByIDTx(ctx, tx, id)
	if err != nil {
		return models.Product{}, err
	}

	if err := tx.Commit(); err != nil {
		return models.Product{}, fmt.Errorf("commit transaction: %w", err)
	}

	return product, nil
}

func (s *ProductService) List(ctx context.Context, includeInactive bool) ([]models.Product, error) {
	query := `
SELECT
	id,
	name,
	base_unit,
	COALESCE(description, ''),
	is_active,
	created_at
FROM product
`
	if !includeInactive {
		query += "WHERE is_active = 1\n"
	}
	query += "ORDER BY id ASC"

	rows, err := s.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("list products: %w", err)
	}
	defer rows.Close()

	result := make([]models.Product, 0)
	for rows.Next() {
		var p models.Product
		if err := rows.Scan(&p.ID, &p.Name, &p.BaseUnit, &p.Description, &p.IsActive, &p.CreatedAt); err != nil {
			return nil, fmt.Errorf("scan product: %w", err)
		}
		result = append(result, p)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate products: %w", err)
	}

	return result, nil
}

func (s *ProductService) Update(ctx context.Context, id int64, name, baseUnit, description string) (models.Product, error) {
	name = strings.TrimSpace(name)
	baseUnit = strings.TrimSpace(baseUnit)

	if name == "" {
		return models.Product{}, errors.New("name is required")
	}
	if !isValidBaseUnit(baseUnit) {
		return models.Product{}, errors.New("base_unit must be one of: g, kg, ml, L, un")
	}

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return models.Product{}, fmt.Errorf("begin transaction: %w", err)
	}
	defer func() {
		_ = tx.Rollback()
	}()

	res, err := tx.ExecContext(
		ctx,
		`UPDATE product SET name = ?, base_unit = ?, description = ? WHERE id = ?`,
		name,
		baseUnit,
		nullableText(description),
		id,
	)
	if err != nil {
		return models.Product{}, fmt.Errorf("update product: %w", err)
	}

	affected, err := res.RowsAffected()
	if err != nil {
		return models.Product{}, fmt.Errorf("read updated rows: %w", err)
	}
	if affected == 0 {
		return models.Product{}, ErrProductNotFound
	}

	product, err := getProductByIDTx(ctx, tx, id)
	if err != nil {
		return models.Product{}, err
	}

	if err := tx.Commit(); err != nil {
		return models.Product{}, fmt.Errorf("commit transaction: %w", err)
	}

	return product, nil
}

func (s *ProductService) SetActive(ctx context.Context, id int64, isActive bool) error {
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

	res, err := tx.ExecContext(ctx, `UPDATE product SET is_active = ? WHERE id = ?`, activeValue, id)
	if err != nil {
		return fmt.Errorf("update product active flag: %w", err)
	}

	affected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("read updated rows: %w", err)
	}
	if affected == 0 {
		return ErrProductNotFound
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit transaction: %w", err)
	}
	return nil
}

func getProductByIDTx(ctx context.Context, tx *sql.Tx, id int64) (models.Product, error) {
	var p models.Product
	err := tx.QueryRowContext(
		ctx,
		`SELECT id, name, base_unit, COALESCE(description, ''), is_active, created_at FROM product WHERE id = ?`,
		id,
	).Scan(&p.ID, &p.Name, &p.BaseUnit, &p.Description, &p.IsActive, &p.CreatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.Product{}, ErrProductNotFound
		}
		return models.Product{}, fmt.Errorf("get product by id: %w", err)
	}
	return p, nil
}

func isValidBaseUnit(baseUnit string) bool {
	return baseUnit == "g" || baseUnit == "kg" || baseUnit == "ml" || baseUnit == "L" || baseUnit == "un"
}
