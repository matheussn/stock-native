package services

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"stock/internal/models"
)

type MovementProductItemInput struct {
	ProductVariationID int64 `json:"product_variation_id"`
	Quantity           int64 `json:"quantity"`
}

type MovementGroupItemInput struct {
	ProductGroupID int64 `json:"product_group_id"`
	Quantity       int64 `json:"quantity"`
}

type CreateMovementInput struct {
	AssistentialWorkID int64                      `json:"assistential_work_id"`
	InstitutionID      int64                      `json:"institution_id"`
	Type               string                     `json:"type"`
	Notes              string                     `json:"notes"`
	ProductItems       []MovementProductItemInput `json:"product_items"`
	GroupItems         []MovementGroupItemInput   `json:"group_items"`
}

type MovementService struct {
	db *sql.DB
}

func NewMovementService(db *sql.DB) *MovementService {
	return &MovementService{db: db}
}

func (s *MovementService) Create(ctx context.Context, input CreateMovementInput) (models.Movement, error) {
	movementType := strings.TrimSpace(input.Type)
	if movementType != "in" && movementType != "out" {
		return models.Movement{}, errors.New("type must be 'in' or 'out'")
	}
	if input.AssistentialWorkID <= 0 {
		return models.Movement{}, errors.New("assistential_work_id must be greater than zero")
	}
	if input.InstitutionID < 0 {
		return models.Movement{}, errors.New("institution_id cannot be negative")
	}
	if len(input.ProductItems) == 0 && len(input.GroupItems) == 0 {
		return models.Movement{}, errors.New("movement must contain at least one product item or group item")
	}

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return models.Movement{}, fmt.Errorf("begin transaction: %w", err)
	}
	defer func() {
		_ = tx.Rollback()
	}()

	storedInstitutionID := int64(0)
	if movementType == "out" && input.InstitutionID > 0 {
		if err := ensureActiveInstitutionTx(ctx, tx, input.InstitutionID); err != nil {
			return models.Movement{}, fmt.Errorf("validate institution: %w", err)
		}
		storedInstitutionID = input.InstitutionID
	}

	res, err := tx.ExecContext(
		ctx,
		`INSERT INTO movement (assistential_work_id, institution_id, type, notes) VALUES (?, ?, ?, ?)`,
		input.AssistentialWorkID,
		nullableInt64(storedInstitutionID),
		movementType,
		nullableText(input.Notes),
	)
	if err != nil {
		return models.Movement{}, fmt.Errorf("insert movement: %w", err)
	}

	movementID, err := res.LastInsertId()
	if err != nil {
		return models.Movement{}, fmt.Errorf("read inserted movement id: %w", err)
	}

	stockDeltaByVariation := make(map[int64]int64)
	deltaSign := int64(1)
	if movementType == "out" {
		deltaSign = -1
	}

	for _, item := range input.ProductItems {
		if item.ProductVariationID <= 0 {
			return models.Movement{}, errors.New("product_variation_id must be greater than zero")
		}
		if item.Quantity <= 0 {
			return models.Movement{}, errors.New("movement product quantity must be greater than zero")
		}

		if _, err := tx.ExecContext(
			ctx,
			`INSERT INTO movement_product_item (movement_id, product_variation_id, quantity) VALUES (?, ?, ?)`,
			movementID,
			item.ProductVariationID,
			item.Quantity,
		); err != nil {
			return models.Movement{}, fmt.Errorf("insert movement product item: %w", err)
		}

		stockDeltaByVariation[item.ProductVariationID] += deltaSign * item.Quantity
	}

	for _, groupItem := range input.GroupItems {
		if groupItem.ProductGroupID <= 0 {
			return models.Movement{}, errors.New("product_group_id must be greater than zero")
		}
		if groupItem.Quantity <= 0 {
			return models.Movement{}, errors.New("movement group quantity must be greater than zero")
		}

		groupRes, err := tx.ExecContext(
			ctx,
			`INSERT INTO movement_group_item (movement_id, product_group_id, quantity) VALUES (?, ?, ?)`,
			movementID,
			groupItem.ProductGroupID,
			groupItem.Quantity,
		)
		if err != nil {
			return models.Movement{}, fmt.Errorf("insert movement group item: %w", err)
		}

		movementGroupItemID, err := groupRes.LastInsertId()
		if err != nil {
			return models.Movement{}, fmt.Errorf("read inserted movement group item id: %w", err)
		}

		resolutions, err := resolveGroupItem(ctx, tx, movementType, groupItem.ProductGroupID, groupItem.Quantity, stockDeltaByVariation)
		if err != nil {
			return models.Movement{}, err
		}

		for _, resolution := range resolutions {
			if _, err := tx.ExecContext(
				ctx,
				`INSERT INTO movement_group_item_resolution (movement_group_item_id, product_variation_id, quantity) VALUES (?, ?, ?)`,
				movementGroupItemID,
				resolution.ProductVariationID,
				resolution.Quantity,
			); err != nil {
				return models.Movement{}, fmt.Errorf("insert movement group item resolution: %w", err)
			}

			stockDeltaByVariation[resolution.ProductVariationID] += deltaSign * resolution.Quantity
		}
	}

	if err := applyStockDeltas(ctx, tx, stockDeltaByVariation); err != nil {
		return models.Movement{}, err
	}

	movement, err := getMovementByIDTx(ctx, tx, movementID)
	if err != nil {
		return models.Movement{}, err
	}

	if err := tx.Commit(); err != nil {
		return models.Movement{}, fmt.Errorf("commit transaction: %w", err)
	}

	return movement, nil
}

func (s *MovementService) List(ctx context.Context, movementType string, assistentialWorkID int64) ([]models.Movement, error) {
	query := `
SELECT
	id,
	assistential_work_id,
	COALESCE(institution_id, 0),
	type,
	COALESCE(notes, ''),
	created_at
FROM movement
WHERE 1 = 1
`
	args := make([]any, 0)

	movementType = strings.TrimSpace(movementType)
	if movementType != "" {
		if movementType != "in" && movementType != "out" {
			return nil, errors.New("type filter must be 'in' or 'out'")
		}
		query += "AND type = ?\n"
		args = append(args, movementType)
	}
	if assistentialWorkID > 0 {
		query += "AND assistential_work_id = ?\n"
		args = append(args, assistentialWorkID)
	}
	query += "ORDER BY created_at DESC, id DESC"

	rows, err := s.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("list movements: %w", err)
	}
	defer rows.Close()

	result := make([]models.Movement, 0)
	for rows.Next() {
		var m models.Movement
		if err := rows.Scan(&m.ID, &m.AssistentialWorkID, &m.InstitutionID, &m.Type, &m.Notes, &m.CreatedAt); err != nil {
			return nil, fmt.Errorf("scan movement: %w", err)
		}
		result = append(result, m)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate movements: %w", err)
	}

	return result, nil
}

func (s *MovementService) ListProductItemsByMovement(ctx context.Context, movementID int64) ([]models.MovementProductItem, error) {
	rows, err := s.db.QueryContext(
		ctx,
		`SELECT id, movement_id, product_variation_id, quantity FROM movement_product_item WHERE movement_id = ? ORDER BY id ASC`,
		movementID,
	)
	if err != nil {
		return nil, fmt.Errorf("list movement product items: %w", err)
	}
	defer rows.Close()

	result := make([]models.MovementProductItem, 0)
	for rows.Next() {
		var item models.MovementProductItem
		if err := rows.Scan(&item.ID, &item.MovementID, &item.ProductVariationID, &item.Quantity); err != nil {
			return nil, fmt.Errorf("scan movement product item: %w", err)
		}
		result = append(result, item)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate movement product items: %w", err)
	}
	return result, nil
}

func (s *MovementService) ListGroupItemsByMovement(ctx context.Context, movementID int64) ([]models.MovementGroupItem, error) {
	rows, err := s.db.QueryContext(
		ctx,
		`SELECT id, movement_id, product_group_id, quantity FROM movement_group_item WHERE movement_id = ? ORDER BY id ASC`,
		movementID,
	)
	if err != nil {
		return nil, fmt.Errorf("list movement group items: %w", err)
	}
	defer rows.Close()

	result := make([]models.MovementGroupItem, 0)
	for rows.Next() {
		var item models.MovementGroupItem
		if err := rows.Scan(&item.ID, &item.MovementID, &item.ProductGroupID, &item.Quantity); err != nil {
			return nil, fmt.Errorf("scan movement group item: %w", err)
		}
		result = append(result, item)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate movement group items: %w", err)
	}
	return result, nil
}

func (s *MovementService) ListGroupItemResolutions(ctx context.Context, movementGroupItemID int64) ([]models.MovementGroupItemResolution, error) {
	rows, err := s.db.QueryContext(
		ctx,
		`SELECT id, movement_group_item_id, product_variation_id, quantity FROM movement_group_item_resolution WHERE movement_group_item_id = ? ORDER BY id ASC`,
		movementGroupItemID,
	)
	if err != nil {
		return nil, fmt.Errorf("list movement group item resolutions: %w", err)
	}
	defer rows.Close()

	result := make([]models.MovementGroupItemResolution, 0)
	for rows.Next() {
		var item models.MovementGroupItemResolution
		if err := rows.Scan(&item.ID, &item.MovementGroupItemID, &item.ProductVariationID, &item.Quantity); err != nil {
			return nil, fmt.Errorf("scan movement group item resolution: %w", err)
		}
		result = append(result, item)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate movement group item resolutions: %w", err)
	}
	return result, nil
}

type groupResolution struct {
	ProductVariationID int64
	Quantity           int64
}

func resolveGroupItem(ctx context.Context, tx *sql.Tx, movementType string, productGroupID, groupQty int64, stockDeltaByVariation map[int64]int64) ([]groupResolution, error) {
	groupRows, err := tx.QueryContext(
		ctx,
		`SELECT product_id, base_quantity FROM product_group_item WHERE product_group_id = ?`,
		productGroupID,
	)
	if err != nil {
		return nil, fmt.Errorf("list product group items for resolution: %w", err)
	}
	defer groupRows.Close()

	resolutions := make([]groupResolution, 0)
	for groupRows.Next() {
		var productID, baseQuantity int64
		if err := groupRows.Scan(&productID, &baseQuantity); err != nil {
			return nil, fmt.Errorf("scan product group item for resolution: %w", err)
		}

		requiredBase := baseQuantity * groupQty
		productResolutions, err := resolveProductRequirement(ctx, tx, movementType, productID, requiredBase, stockDeltaByVariation)
		if err != nil {
			return nil, err
		}
		resolutions = append(resolutions, productResolutions...)
	}
	if err := groupRows.Err(); err != nil {
		return nil, fmt.Errorf("iterate product group items for resolution: %w", err)
	}

	if len(resolutions) == 0 {
		return nil, fmt.Errorf("product group %d has no items to resolve", productGroupID)
	}

	return resolutions, nil
}

func resolveProductRequirement(ctx context.Context, tx *sql.Tx, movementType string, productID, requiredBase int64, stockDeltaByVariation map[int64]int64) ([]groupResolution, error) {
	type candidate struct {
		variationID int64
		baseQty     int64
		stock       int64
	}

	rows, err := tx.QueryContext(
		ctx,
		`
SELECT
	id,
	base_quantity,
	current_stock
FROM product_variation
WHERE product_id = ? AND is_active = 1
ORDER BY base_quantity DESC, id ASC
`,
		productID,
	)
	if err != nil {
		return nil, fmt.Errorf("list product variations for resolution: %w", err)
	}
	defer rows.Close()

	candidates := make([]candidate, 0)
	for rows.Next() {
		var c candidate
		if err := rows.Scan(&c.variationID, &c.baseQty, &c.stock); err != nil {
			return nil, fmt.Errorf("scan product variation for resolution: %w", err)
		}
		candidates = append(candidates, c)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate product variations for resolution: %w", err)
	}

	if len(candidates) == 0 {
		return nil, fmt.Errorf("no active variations available for product %d", productID)
	}

	remaining := requiredBase
	result := make([]groupResolution, 0)

	for _, c := range candidates {
		if c.baseQty <= 0 {
			continue
		}

		maxUnits := remaining / c.baseQty
		if maxUnits <= 0 {
			continue
		}

		allocUnits := maxUnits
		if movementType == "out" {
			effectiveStock := c.stock + stockDeltaByVariation[c.variationID]
			if effectiveStock < 0 {
				effectiveStock = 0
			}
			if allocUnits > effectiveStock {
				allocUnits = effectiveStock
			}
		}
		if allocUnits <= 0 {
			continue
		}

		result = append(result, groupResolution{
			ProductVariationID: c.variationID,
			Quantity:           allocUnits,
		})
		remaining -= allocUnits * c.baseQty
		if remaining == 0 {
			break
		}
	}

	if remaining > 0 {
		if movementType == "out" {
			return nil, fmt.Errorf("insufficient stock to resolve product %d in product group", productID)
		}
		return nil, fmt.Errorf("cannot resolve product %d with available variations", productID)
	}

	return result, nil
}

func applyStockDeltas(ctx context.Context, tx *sql.Tx, deltas map[int64]int64) error {
	for variationID, delta := range deltas {
		if delta == 0 {
			continue
		}

		res, err := tx.ExecContext(
			ctx,
			`UPDATE product_variation SET current_stock = current_stock + ? WHERE id = ? AND current_stock + ? >= 0`,
			delta,
			variationID,
			delta,
		)
		if err != nil {
			return fmt.Errorf("update stock for variation %d: %w", variationID, err)
		}

		affected, err := res.RowsAffected()
		if err != nil {
			return fmt.Errorf("read updated stock rows for variation %d: %w", variationID, err)
		}
		if affected == 0 {
			return fmt.Errorf("insufficient stock for product variation %d", variationID)
		}
	}
	return nil
}

func getMovementByIDTx(ctx context.Context, tx *sql.Tx, movementID int64) (models.Movement, error) {
	var m models.Movement
	err := tx.QueryRowContext(
		ctx,
		`SELECT id, assistential_work_id, COALESCE(institution_id, 0), type, COALESCE(notes, ''), created_at FROM movement WHERE id = ?`,
		movementID,
	).Scan(&m.ID, &m.AssistentialWorkID, &m.InstitutionID, &m.Type, &m.Notes, &m.CreatedAt)
	if err != nil {
		return models.Movement{}, fmt.Errorf("get movement by id: %w", err)
	}
	return m, nil
}

func nullableInt64(value int64) any {
	if value <= 0 {
		return nil
	}
	return value
}
