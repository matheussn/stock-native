package services

import (
	"context"
	"database/sql"
	"fmt"
)

type VariationStockStatus struct {
	ProductID            int64  `json:"product_id"`
	ProductName          string `json:"product_name"`
	BaseUnit             string `json:"base_unit"`
	VariationID          int64  `json:"variation_id"`
	VariationDescription string `json:"variation_description"`
	BaseQuantity         int64  `json:"base_quantity"`
	CurrentStock         int64  `json:"current_stock"`
}

type MonthlyDemandItem struct {
	ProductID            int64  `json:"product_id"`
	ProductName          string `json:"product_name"`
	BaseUnit             string `json:"base_unit"`
	RequiredBaseQuantity int64  `json:"required_base_quantity"`
}

type StockCoverageItem struct {
	ProductID             int64  `json:"product_id"`
	ProductName           string `json:"product_name"`
	BaseUnit              string `json:"base_unit"`
	RequiredBaseQuantity  int64  `json:"required_base_quantity"`
	AvailableBaseQuantity int64  `json:"available_base_quantity"`
	ShortfallBaseQuantity int64  `json:"shortfall_base_quantity"`
	IsCovered             bool   `json:"is_covered"`
}

type AssistentialWorkOutflowKgItem struct {
	AssistentialWorkID   int64   `json:"assistential_work_id"`
	AssistentialWorkName string  `json:"assistential_work_name"`
	OutputKg             float64 `json:"output_kg"`
}

type MonthlyMovementFlowKgItem struct {
	MonthKey string  `json:"month_key"`
	InputKg  float64 `json:"input_kg"`
	OutputKg float64 `json:"output_kg"`
}

type StockService struct {
	db *sql.DB
}

func NewStockService(db *sql.DB) *StockService {
	return &StockService{db: db}
}

func (s *StockService) ListVariationStatus(ctx context.Context, includeInactive bool) ([]VariationStockStatus, error) {
	query := `
SELECT
	p.id,
	p.name,
	p.base_unit,
	pv.id,
	COALESCE(pv.description, ''),
	pv.base_quantity,
	pv.current_stock
FROM product_variation pv
JOIN product p ON p.id = pv.product_id
WHERE p.is_active = 1
`
	if !includeInactive {
		query += "AND pv.is_active = 1\n"
	}
	query += "ORDER BY p.name ASC, pv.base_quantity DESC, pv.id ASC"

	rows, err := s.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("list variation stock status: %w", err)
	}
	defer rows.Close()

	result := make([]VariationStockStatus, 0)
	for rows.Next() {
		var item VariationStockStatus
		if err := rows.Scan(
			&item.ProductID,
			&item.ProductName,
			&item.BaseUnit,
			&item.VariationID,
			&item.VariationDescription,
			&item.BaseQuantity,
			&item.CurrentStock,
		); err != nil {
			return nil, fmt.Errorf("scan variation stock status: %w", err)
		}
		result = append(result, item)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate variation stock status: %w", err)
	}

	return result, nil
}

func (s *StockService) GetMonthlyDemandProjection(ctx context.Context) ([]MonthlyDemandItem, error) {
	rows, err := s.db.QueryContext(
		ctx,
		`
SELECT
	p.id,
	p.name,
	p.base_unit,
	COALESCE(SUM(pgi.base_quantity), 0) AS required_base_quantity
FROM family_group_assignment fga
JOIN family f ON f.id = fga.family_id
JOIN product_group pg ON pg.id = fga.product_group_id
JOIN product_group_item pgi ON pgi.product_group_id = pg.id
JOIN product p ON p.id = pgi.product_id
WHERE fga.ended_at IS NULL
  AND f.is_active = 1
  AND pg.is_active = 1
  AND p.is_active = 1
GROUP BY p.id, p.name, p.base_unit
ORDER BY p.name ASC
`,
	)
	if err != nil {
		return nil, fmt.Errorf("query monthly demand projection: %w", err)
	}
	defer rows.Close()

	result := make([]MonthlyDemandItem, 0)
	for rows.Next() {
		var item MonthlyDemandItem
		if err := rows.Scan(&item.ProductID, &item.ProductName, &item.BaseUnit, &item.RequiredBaseQuantity); err != nil {
			return nil, fmt.Errorf("scan monthly demand projection: %w", err)
		}
		result = append(result, item)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate monthly demand projection: %w", err)
	}
	return result, nil
}

func (s *StockService) GetCoverageCheck(ctx context.Context) ([]StockCoverageItem, error) {
	demand, err := s.GetMonthlyDemandProjection(ctx)
	if err != nil {
		return nil, err
	}
	if len(demand) == 0 {
		return []StockCoverageItem{}, nil
	}

	query := `
SELECT COALESCE(SUM(pv.current_stock * pv.base_quantity), 0)
FROM product_variation pv
WHERE pv.product_id = ? AND pv.is_active = 1
`

	result := make([]StockCoverageItem, 0, len(demand))
	for _, d := range demand {
		var availableBase int64
		if err := s.db.QueryRowContext(ctx, query, d.ProductID).Scan(&availableBase); err != nil {
			return nil, fmt.Errorf("query available base quantity for product %d: %w", d.ProductID, err)
		}

		shortfall := int64(0)
		if availableBase < d.RequiredBaseQuantity {
			shortfall = d.RequiredBaseQuantity - availableBase
		}

		result = append(result, StockCoverageItem{
			ProductID:             d.ProductID,
			ProductName:           d.ProductName,
			BaseUnit:              d.BaseUnit,
			RequiredBaseQuantity:  d.RequiredBaseQuantity,
			AvailableBaseQuantity: availableBase,
			ShortfallBaseQuantity: shortfall,
			IsCovered:             shortfall == 0,
		})
	}

	return result, nil
}

func (s *StockService) GetAssistentialWorkOutflowKg(ctx context.Context) ([]AssistentialWorkOutflowKgItem, error) {
	rows, err := s.db.QueryContext(
		ctx,
		`
WITH outflow AS (
	SELECT
		m.assistential_work_id,
		SUM(
			CASE
				WHEN p.base_unit = 'kg' THEN CAST(mpi.quantity * pv.base_quantity AS REAL)
				WHEN p.base_unit = 'g' THEN CAST(mpi.quantity * pv.base_quantity AS REAL) / 1000.0
				ELSE 0
			END
		) AS total_kg
	FROM movement m
	JOIN movement_product_item mpi ON mpi.movement_id = m.id
	JOIN product_variation pv ON pv.id = mpi.product_variation_id
	JOIN product p ON p.id = pv.product_id
	WHERE m.type = 'out'
	  AND p.base_unit IN ('g', 'kg')
	GROUP BY m.assistential_work_id

	UNION ALL

	SELECT
		m.assistential_work_id,
		SUM(
			CASE
				WHEN p.base_unit = 'kg' THEN CAST(mgir.quantity * pv.base_quantity AS REAL)
				WHEN p.base_unit = 'g' THEN CAST(mgir.quantity * pv.base_quantity AS REAL) / 1000.0
				ELSE 0
			END
		) AS total_kg
	FROM movement m
	JOIN movement_group_item mgi ON mgi.movement_id = m.id
	JOIN movement_group_item_resolution mgir ON mgir.movement_group_item_id = mgi.id
	JOIN product_variation pv ON pv.id = mgir.product_variation_id
	JOIN product p ON p.id = pv.product_id
	WHERE m.type = 'out'
	  AND p.base_unit IN ('g', 'kg')
	GROUP BY m.assistential_work_id
)
SELECT
	aw.id,
	aw.name,
	ROUND(COALESCE(SUM(outflow.total_kg), 0), 3) AS total_kg
FROM assistential_work aw
JOIN outflow ON outflow.assistential_work_id = aw.id
GROUP BY aw.id, aw.name
HAVING COALESCE(SUM(outflow.total_kg), 0) > 0
ORDER BY total_kg DESC, aw.name ASC
`,
	)
	if err != nil {
		return nil, fmt.Errorf("query assistential work outflow kg: %w", err)
	}
	defer rows.Close()

	result := make([]AssistentialWorkOutflowKgItem, 0)
	for rows.Next() {
		var item AssistentialWorkOutflowKgItem
		if err := rows.Scan(&item.AssistentialWorkID, &item.AssistentialWorkName, &item.OutputKg); err != nil {
			return nil, fmt.Errorf("scan assistential work outflow kg: %w", err)
		}
		result = append(result, item)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate assistential work outflow kg: %w", err)
	}

	return result, nil
}

func (s *StockService) GetMonthlyMovementFlowKg(ctx context.Context) ([]MonthlyMovementFlowKgItem, error) {
	rows, err := s.db.QueryContext(
		ctx,
		`
WITH monthly_flow AS (
	SELECT
		strftime('%Y-%m', m.created_at) AS month_key,
		m.type,
		SUM(
			CASE
				WHEN p.base_unit = 'kg' THEN CAST(mpi.quantity * pv.base_quantity AS REAL)
				WHEN p.base_unit = 'g' THEN CAST(mpi.quantity * pv.base_quantity AS REAL) / 1000.0
				ELSE 0
			END
		) AS total_kg
	FROM movement m
	JOIN movement_product_item mpi ON mpi.movement_id = m.id
	JOIN product_variation pv ON pv.id = mpi.product_variation_id
	JOIN product p ON p.id = pv.product_id
	WHERE p.base_unit IN ('g', 'kg')
	GROUP BY month_key, m.type

	UNION ALL

	SELECT
		strftime('%Y-%m', m.created_at) AS month_key,
		m.type,
		SUM(
			CASE
				WHEN p.base_unit = 'kg' THEN CAST(mgir.quantity * pv.base_quantity AS REAL)
				WHEN p.base_unit = 'g' THEN CAST(mgir.quantity * pv.base_quantity AS REAL) / 1000.0
				ELSE 0
			END
		) AS total_kg
	FROM movement m
	JOIN movement_group_item mgi ON mgi.movement_id = m.id
	JOIN movement_group_item_resolution mgir ON mgir.movement_group_item_id = mgi.id
	JOIN product_variation pv ON pv.id = mgir.product_variation_id
	JOIN product p ON p.id = pv.product_id
	WHERE p.base_unit IN ('g', 'kg')
	GROUP BY month_key, m.type
)
SELECT
	month_key,
	ROUND(SUM(CASE WHEN type = 'in' THEN total_kg ELSE 0 END), 3) AS input_kg,
	ROUND(SUM(CASE WHEN type = 'out' THEN total_kg ELSE 0 END), 3) AS output_kg
FROM monthly_flow
GROUP BY month_key
ORDER BY month_key ASC
`,
	)
	if err != nil {
		return nil, fmt.Errorf("query monthly movement flow kg: %w", err)
	}
	defer rows.Close()

	result := make([]MonthlyMovementFlowKgItem, 0)
	for rows.Next() {
		var item MonthlyMovementFlowKgItem
		if err := rows.Scan(&item.MonthKey, &item.InputKg, &item.OutputKg); err != nil {
			return nil, fmt.Errorf("scan monthly movement flow kg: %w", err)
		}
		result = append(result, item)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate monthly movement flow kg: %w", err)
	}

	return result, nil
}
