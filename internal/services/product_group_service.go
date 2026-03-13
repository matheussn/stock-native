package services

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"stock/internal/models"
)

var ErrProductGroupNotFound = errors.New("product group not found")
var ErrProductGroupItemNotFound = errors.New("product group item not found")

type ProductGroupService struct {
	db *sql.DB
}

func NewProductGroupService(db *sql.DB) *ProductGroupService {
	return &ProductGroupService{db: db}
}

func (s *ProductGroupService) Create(ctx context.Context, name, description string) (models.ProductGroup, error) {
	name = strings.TrimSpace(name)
	if name == "" {
		return models.ProductGroup{}, errors.New("name is required")
	}

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return models.ProductGroup{}, fmt.Errorf("begin transaction: %w", err)
	}
	defer func() {
		_ = tx.Rollback()
	}()

	res, err := tx.ExecContext(
		ctx,
		`INSERT INTO product_group (name, description, is_active) VALUES (?, ?, 1)`,
		name,
		nullableText(description),
	)
	if err != nil {
		return models.ProductGroup{}, fmt.Errorf("insert product group: %w", err)
	}

	id, err := res.LastInsertId()
	if err != nil {
		return models.ProductGroup{}, fmt.Errorf("read inserted id: %w", err)
	}

	group, err := getProductGroupByIDTx(ctx, tx, id)
	if err != nil {
		return models.ProductGroup{}, err
	}

	if err := tx.Commit(); err != nil {
		return models.ProductGroup{}, fmt.Errorf("commit transaction: %w", err)
	}

	return group, nil
}

func (s *ProductGroupService) List(ctx context.Context, includeInactive bool) ([]models.ProductGroup, error) {
	query := `
SELECT
	id,
	name,
	COALESCE(description, ''),
	is_active,
	created_at
FROM product_group
`
	if !includeInactive {
		query += "WHERE is_active = 1\n"
	}
	query += "ORDER BY id ASC"

	rows, err := s.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("list product groups: %w", err)
	}
	defer rows.Close()

	result := make([]models.ProductGroup, 0)
	for rows.Next() {
		var g models.ProductGroup
		if err := rows.Scan(&g.ID, &g.Name, &g.Description, &g.IsActive, &g.CreatedAt); err != nil {
			return nil, fmt.Errorf("scan product group: %w", err)
		}
		result = append(result, g)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate product groups: %w", err)
	}

	return result, nil
}

func (s *ProductGroupService) Update(ctx context.Context, id int64, name, description string) (models.ProductGroup, error) {
	name = strings.TrimSpace(name)
	if name == "" {
		return models.ProductGroup{}, errors.New("name is required")
	}

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return models.ProductGroup{}, fmt.Errorf("begin transaction: %w", err)
	}
	defer func() {
		_ = tx.Rollback()
	}()

	res, err := tx.ExecContext(
		ctx,
		`UPDATE product_group SET name = ?, description = ? WHERE id = ?`,
		name,
		nullableText(description),
		id,
	)
	if err != nil {
		return models.ProductGroup{}, fmt.Errorf("update product group: %w", err)
	}

	affected, err := res.RowsAffected()
	if err != nil {
		return models.ProductGroup{}, fmt.Errorf("read updated rows: %w", err)
	}
	if affected == 0 {
		return models.ProductGroup{}, ErrProductGroupNotFound
	}

	group, err := getProductGroupByIDTx(ctx, tx, id)
	if err != nil {
		return models.ProductGroup{}, err
	}

	if err := tx.Commit(); err != nil {
		return models.ProductGroup{}, fmt.Errorf("commit transaction: %w", err)
	}

	return group, nil
}

func (s *ProductGroupService) SetActive(ctx context.Context, id int64, isActive bool) error {
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

	res, err := tx.ExecContext(ctx, `UPDATE product_group SET is_active = ? WHERE id = ?`, activeValue, id)
	if err != nil {
		return fmt.Errorf("update product group active flag: %w", err)
	}

	affected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("read updated rows: %w", err)
	}
	if affected == 0 {
		return ErrProductGroupNotFound
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit transaction: %w", err)
	}
	return nil
}

func (s *ProductGroupService) ListItems(ctx context.Context, productGroupID int64) ([]models.ProductGroupItem, error) {
	rows, err := s.db.QueryContext(
		ctx,
		`SELECT id, product_group_id, product_id, base_quantity FROM product_group_item WHERE product_group_id = ? ORDER BY id ASC`,
		productGroupID,
	)
	if err != nil {
		return nil, fmt.Errorf("list product group items: %w", err)
	}
	defer rows.Close()

	result := make([]models.ProductGroupItem, 0)
	for rows.Next() {
		var item models.ProductGroupItem
		if err := rows.Scan(&item.ID, &item.ProductGroupID, &item.ProductID, &item.BaseQuantity); err != nil {
			return nil, fmt.Errorf("scan product group item: %w", err)
		}
		result = append(result, item)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate product group items: %w", err)
	}

	return result, nil
}

func (s *ProductGroupService) UpsertItem(ctx context.Context, productGroupID, productID, baseQuantity int64) (models.ProductGroupItem, error) {
	if productGroupID <= 0 {
		return models.ProductGroupItem{}, errors.New("product_group_id must be greater than zero")
	}
	if productID <= 0 {
		return models.ProductGroupItem{}, errors.New("product_id must be greater than zero")
	}
	if baseQuantity <= 0 {
		return models.ProductGroupItem{}, errors.New("base_quantity must be greater than zero")
	}

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return models.ProductGroupItem{}, fmt.Errorf("begin transaction: %w", err)
	}
	defer func() {
		_ = tx.Rollback()
	}()

	_, err = tx.ExecContext(
		ctx,
		`
INSERT INTO product_group_item (product_group_id, product_id, base_quantity)
VALUES (?, ?, ?)
ON CONFLICT(product_group_id, product_id)
DO UPDATE SET base_quantity = excluded.base_quantity
`,
		productGroupID,
		productID,
		baseQuantity,
	)
	if err != nil {
		return models.ProductGroupItem{}, fmt.Errorf("upsert product group item: %w", err)
	}

	item, err := getProductGroupItemByKeysTx(ctx, tx, productGroupID, productID)
	if err != nil {
		return models.ProductGroupItem{}, err
	}

	if err := tx.Commit(); err != nil {
		return models.ProductGroupItem{}, fmt.Errorf("commit transaction: %w", err)
	}
	return item, nil
}

func (s *ProductGroupService) RemoveItem(ctx context.Context, productGroupID, productID int64) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("begin transaction: %w", err)
	}
	defer func() {
		_ = tx.Rollback()
	}()

	res, err := tx.ExecContext(
		ctx,
		`DELETE FROM product_group_item WHERE product_group_id = ? AND product_id = ?`,
		productGroupID,
		productID,
	)
	if err != nil {
		return fmt.Errorf("remove product group item: %w", err)
	}

	affected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("read removed rows: %w", err)
	}
	if affected == 0 {
		return ErrProductGroupItemNotFound
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit transaction: %w", err)
	}
	return nil
}

func getProductGroupByIDTx(ctx context.Context, tx *sql.Tx, id int64) (models.ProductGroup, error) {
	var g models.ProductGroup
	err := tx.QueryRowContext(
		ctx,
		`SELECT id, name, COALESCE(description, ''), is_active, created_at FROM product_group WHERE id = ?`,
		id,
	).Scan(&g.ID, &g.Name, &g.Description, &g.IsActive, &g.CreatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.ProductGroup{}, ErrProductGroupNotFound
		}
		return models.ProductGroup{}, fmt.Errorf("get product group by id: %w", err)
	}
	return g, nil
}

func getProductGroupItemByIDTx(ctx context.Context, tx *sql.Tx, id int64) (models.ProductGroupItem, error) {
	var item models.ProductGroupItem
	err := tx.QueryRowContext(
		ctx,
		`SELECT id, product_group_id, product_id, base_quantity FROM product_group_item WHERE id = ?`,
		id,
	).Scan(&item.ID, &item.ProductGroupID, &item.ProductID, &item.BaseQuantity)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.ProductGroupItem{}, ErrProductGroupItemNotFound
		}
		return models.ProductGroupItem{}, fmt.Errorf("get product group item by id: %w", err)
	}
	return item, nil
}

func getProductGroupItemByKeysTx(ctx context.Context, tx *sql.Tx, productGroupID, productID int64) (models.ProductGroupItem, error) {
	var item models.ProductGroupItem
	err := tx.QueryRowContext(
		ctx,
		`SELECT id, product_group_id, product_id, base_quantity FROM product_group_item WHERE product_group_id = ? AND product_id = ?`,
		productGroupID,
		productID,
	).Scan(&item.ID, &item.ProductGroupID, &item.ProductID, &item.BaseQuantity)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.ProductGroupItem{}, ErrProductGroupItemNotFound
		}
		return models.ProductGroupItem{}, fmt.Errorf("get product group item by keys: %w", err)
	}
	return item, nil
}
