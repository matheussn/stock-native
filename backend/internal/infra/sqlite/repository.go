package sqlite

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"stock/backend/internal/application"
	"stock/backend/internal/domain"
)

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) ListCategories(ctx context.Context) ([]domain.Category, error) {
	rows, err := r.db.QueryContext(ctx, `SELECT id, name, active, created_at, updated_at FROM categories ORDER BY name`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := make([]domain.Category, 0)
	for rows.Next() {
		var item domain.Category
		if err := rows.Scan(&item.ID, &item.Name, &item.Active, &item.CreatedAt, &item.UpdatedAt); err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return items, rows.Err()
}

func (r *Repository) CreateCategory(ctx context.Context, item domain.Category) error {
	_, err := r.db.ExecContext(ctx, `INSERT INTO categories(id, name, active, created_at, updated_at) VALUES (?, ?, ?, ?, ?)`, item.ID, item.Name, item.Active, item.CreatedAt, item.UpdatedAt)
	return err
}

func (r *Repository) UpdateCategory(ctx context.Context, item domain.Category) error {
	result, err := r.db.ExecContext(ctx, `UPDATE categories SET name = ?, updated_at = ? WHERE id = ?`, item.Name, item.UpdatedAt, item.ID)
	if err != nil {
		return err
	}
	return requireRowsAffected(result, "category not found")
}

func (r *Repository) DeactivateCategory(ctx context.Context, id string) error {
	result, err := r.db.ExecContext(ctx, `UPDATE categories SET active = 0, updated_at = ? WHERE id = ?`, time.Now(), id)
	if err != nil {
		return err
	}
	return requireRowsAffected(result, "category not found")
}

func (r *Repository) ListProducts(ctx context.Context) ([]domain.Product, error) {
	rows, err := r.db.QueryContext(ctx, `SELECT id, name, category_id, measure_unit, package_amount, low_stock_threshold, active, created_at, updated_at FROM products ORDER BY name`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := make([]domain.Product, 0)
	for rows.Next() {
		var item domain.Product
		if err := rows.Scan(&item.ID, &item.Name, &item.CategoryID, &item.MeasureUnit, &item.PackageAmount, &item.LowStockThreshold, &item.Active, &item.CreatedAt, &item.UpdatedAt); err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return items, rows.Err()
}

func (r *Repository) CreateProduct(ctx context.Context, item domain.Product) error {
	_, err := r.db.ExecContext(ctx, `INSERT INTO products(id, name, category_id, unit, measure_unit, package_amount, low_stock_threshold, active, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`, item.ID, item.Name, item.CategoryID, item.MeasureUnit, item.MeasureUnit, item.PackageAmount, item.LowStockThreshold, item.Active, item.CreatedAt, item.UpdatedAt)
	if isForeignKeyError(err) {
		return domain.NewAppError(domain.ErrorValidation, "categoryId is invalid", nil)
	}
	return err
}

func (r *Repository) UpdateProduct(ctx context.Context, item domain.Product) error {
	result, err := r.db.ExecContext(ctx, `UPDATE products SET name = ?, category_id = ?, unit = ?, measure_unit = ?, package_amount = ?, low_stock_threshold = ?, updated_at = ? WHERE id = ?`, item.Name, item.CategoryID, item.MeasureUnit, item.MeasureUnit, item.PackageAmount, item.LowStockThreshold, item.UpdatedAt, item.ID)
	if isForeignKeyError(err) {
		return domain.NewAppError(domain.ErrorValidation, "categoryId is invalid", nil)
	}
	if err != nil {
		return err
	}
	return requireRowsAffected(result, "product not found")
}

func (r *Repository) DeactivateProduct(ctx context.Context, id string) error {
	result, err := r.db.ExecContext(ctx, `UPDATE products SET active = 0, updated_at = ? WHERE id = ?`, time.Now(), id)
	if err != nil {
		return err
	}
	return requireRowsAffected(result, "product not found")
}

func (r *Repository) ListOrigins(ctx context.Context) ([]domain.Origin, error) {
	rows, err := r.db.QueryContext(ctx, `SELECT id, name, active, created_at, updated_at FROM origins ORDER BY name`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := make([]domain.Origin, 0)
	for rows.Next() {
		var item domain.Origin
		if err := rows.Scan(&item.ID, &item.Name, &item.Active, &item.CreatedAt, &item.UpdatedAt); err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return items, rows.Err()
}

func (r *Repository) CreateOrigin(ctx context.Context, item domain.Origin) error {
	_, err := r.db.ExecContext(ctx, `INSERT INTO origins(id, name, active, created_at, updated_at) VALUES (?, ?, ?, ?, ?)`, item.ID, item.Name, item.Active, item.CreatedAt, item.UpdatedAt)
	return err
}

func (r *Repository) UpdateOrigin(ctx context.Context, item domain.Origin) error {
	result, err := r.db.ExecContext(ctx, `UPDATE origins SET name = ?, updated_at = ? WHERE id = ?`, item.Name, item.UpdatedAt, item.ID)
	if err != nil {
		return err
	}
	return requireRowsAffected(result, "origin not found")
}

func (r *Repository) DeactivateOrigin(ctx context.Context, id string) error {
	result, err := r.db.ExecContext(ctx, `UPDATE origins SET active = 0, updated_at = ? WHERE id = ?`, time.Now(), id)
	if err != nil {
		return err
	}
	return requireRowsAffected(result, "origin not found")
}

func (r *Repository) ListDestinations(ctx context.Context) ([]domain.Destination, error) {
	rows, err := r.db.QueryContext(ctx, `SELECT id, name, active, created_at, updated_at FROM destinations ORDER BY name`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := make([]domain.Destination, 0)
	for rows.Next() {
		var item domain.Destination
		if err := rows.Scan(&item.ID, &item.Name, &item.Active, &item.CreatedAt, &item.UpdatedAt); err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return items, rows.Err()
}

func (r *Repository) CreateDestination(ctx context.Context, item domain.Destination) error {
	_, err := r.db.ExecContext(ctx, `INSERT INTO destinations(id, name, active, created_at, updated_at) VALUES (?, ?, ?, ?, ?)`, item.ID, item.Name, item.Active, item.CreatedAt, item.UpdatedAt)
	return err
}

func (r *Repository) UpdateDestination(ctx context.Context, item domain.Destination) error {
	result, err := r.db.ExecContext(ctx, `UPDATE destinations SET name = ?, updated_at = ? WHERE id = ?`, item.Name, item.UpdatedAt, item.ID)
	if err != nil {
		return err
	}
	return requireRowsAffected(result, "destination not found")
}

func (r *Repository) DeactivateDestination(ctx context.Context, id string) error {
	result, err := r.db.ExecContext(ctx, `UPDATE destinations SET active = 0, updated_at = ? WHERE id = ?`, time.Now(), id)
	if err != nil {
		return err
	}
	return requireRowsAffected(result, "destination not found")
}

func (r *Repository) CreateMovementEntry(ctx context.Context, item domain.Movement) error {
	_, err := r.db.ExecContext(ctx, `INSERT INTO movements(id, type, product_id, quantity, movement_date, origin_id, destination_id, note, created_at) VALUES (?, ?, ?, ?, ?, ?, NULL, ?, ?)`, item.ID, item.Type, item.ProductID, item.Quantity, item.MovementDate, item.OriginID, item.Note, item.CreatedAt)
	if isForeignKeyError(err) {
		return domain.NewAppError(domain.ErrorValidation, "invalid product or origin", nil)
	}
	return err
}

func (r *Repository) CreateMovementExit(ctx context.Context, item domain.Movement) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	var currentStock float64
	query := `SELECT COALESCE(SUM(CASE WHEN type = 'entrada' THEN quantity ELSE -quantity END), 0) FROM movements WHERE product_id = ?`
	if err := tx.QueryRowContext(ctx, query, item.ProductID).Scan(&currentStock); err != nil {
		return err
	}
	if currentStock < item.Quantity {
		return domain.NewAppError(domain.ErrorInsufficientStock, "insufficient stock for exit movement", map[string]float64{"currentStock": domain.Round3(currentStock), "requested": item.Quantity})
	}

	_, err = tx.ExecContext(ctx, `INSERT INTO movements(id, type, product_id, quantity, movement_date, origin_id, destination_id, note, created_at) VALUES (?, ?, ?, ?, ?, NULL, ?, ?, ?)`, item.ID, item.Type, item.ProductID, item.Quantity, item.MovementDate, item.DestinationID, item.Note, item.CreatedAt)
	if isForeignKeyError(err) {
		return domain.NewAppError(domain.ErrorValidation, "invalid product or destination", nil)
	}
	if err != nil {
		return err
	}
	return tx.Commit()
}

func (r *Repository) CreateMovementEntriesBatch(ctx context.Context, items []domain.Movement) error {
	if len(items) == 0 {
		return domain.NewAppError(domain.ErrorValidation, "items are required", nil)
	}

	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	for _, item := range items {
		if err := insertMovementEntryTx(ctx, tx, item); err != nil {
			return err
		}
	}

	return tx.Commit()
}

func (r *Repository) CreateMovementExitsBatch(ctx context.Context, items []domain.Movement) error {
	if len(items) == 0 {
		return domain.NewAppError(domain.ErrorValidation, "items are required", nil)
	}

	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	requestedByProduct := map[string]float64{}
	for _, item := range items {
		requestedByProduct[item.ProductID] += item.Quantity
	}

	for productID, requested := range requestedByProduct {
		var currentStock float64
		query := `SELECT COALESCE(SUM(CASE WHEN type = 'entrada' THEN quantity ELSE -quantity END), 0) FROM movements WHERE product_id = ?`
		if err := tx.QueryRowContext(ctx, query, productID).Scan(&currentStock); err != nil {
			return err
		}
		if currentStock < requested {
			return domain.NewAppError(domain.ErrorInsufficientStock, "insufficient stock for exit movement", map[string]float64{
				"currentStock": domain.Round3(currentStock),
				"requested":    domain.Round3(requested),
			})
		}
	}

	for _, item := range items {
		if err := insertMovementExitTx(ctx, tx, item); err != nil {
			return err
		}
	}

	return tx.Commit()
}

func (r *Repository) ListMovements(ctx context.Context, filter application.ListMovementsFilter) ([]domain.Movement, error) {
	args := make([]any, 0)
	where := make([]string, 0)
	if strings.TrimSpace(filter.Type) != "" {
		where = append(where, "m.type = ?")
		args = append(args, filter.Type)
	}
	if strings.TrimSpace(filter.ProductID) != "" {
		where = append(where, "m.product_id = ?")
		args = append(args, filter.ProductID)
	}
	if filter.DateFrom != nil {
		where = append(where, "m.movement_date >= ?")
		args = append(args, *filter.DateFrom)
	}
	if filter.DateTo != nil {
		where = append(where, "m.movement_date <= ?")
		args = append(args, *filter.DateTo)
	}

	query := `SELECT m.id, m.type, m.product_id, p.name, m.quantity, m.movement_date, m.origin_id, o.name, m.destination_id, d.name, m.note, m.created_at
		FROM movements m
		JOIN products p ON p.id = m.product_id
		LEFT JOIN origins o ON o.id = m.origin_id
		LEFT JOIN destinations d ON d.id = m.destination_id`
	if len(where) > 0 {
		query += " WHERE " + strings.Join(where, " AND ")
	}
	query += " ORDER BY m.movement_date DESC, m.created_at DESC"

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := make([]domain.Movement, 0)
	for rows.Next() {
		var item domain.Movement
		var originID sql.NullString
		var originName sql.NullString
		var destinationID sql.NullString
		var destinationName sql.NullString

		if err := rows.Scan(&item.ID, &item.Type, &item.ProductID, &item.ProductName, &item.Quantity, &item.MovementDate, &originID, &originName, &destinationID, &destinationName, &item.Note, &item.CreatedAt); err != nil {
			return nil, err
		}
		if originID.Valid {
			item.OriginID = &originID.String
		}
		if originName.Valid {
			item.OriginName = &originName.String
		}
		if destinationID.Valid {
			item.DestinationID = &destinationID.String
		}
		if destinationName.Valid {
			item.DestinationName = &destinationName.String
		}
		item.Quantity = domain.Round3(item.Quantity)
		items = append(items, item)
	}
	return items, rows.Err()
}

func (r *Repository) GetCurrentStock(ctx context.Context, filter application.CurrentStockFilter) ([]domain.StockItem, error) {
	args := make([]any, 0)
	where := []string{"p.active = 1"}
	if strings.TrimSpace(filter.CategoryID) != "" {
		where = append(where, "p.category_id = ?")
		args = append(args, filter.CategoryID)
	}
	if strings.TrimSpace(filter.Search) != "" {
		where = append(where, "LOWER(p.name) LIKE ?")
		args = append(args, "%"+strings.ToLower(strings.TrimSpace(filter.Search))+"%")
	}

	query := `SELECT p.id, p.name, c.id, c.name, p.measure_unit, p.package_amount, p.low_stock_threshold,
		COALESCE(SUM(CASE WHEN m.type='entrada' THEN m.quantity ELSE -m.quantity END), 0) AS current_stock
		FROM products p
		JOIN categories c ON c.id = p.category_id
		LEFT JOIN movements m ON m.product_id = p.id`
	if len(where) > 0 {
		query += " WHERE " + strings.Join(where, " AND ")
	}
	query += " GROUP BY p.id, p.name, c.id, c.name, p.measure_unit, p.package_amount, p.low_stock_threshold ORDER BY p.name"

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := make([]domain.StockItem, 0)
	for rows.Next() {
		var item domain.StockItem
		if err := rows.Scan(&item.ProductID, &item.ProductName, &item.CategoryID, &item.CategoryName, &item.MeasureUnit, &item.PackageAmount, &item.LowStockThreshold, &item.CurrentStock); err != nil {
			return nil, err
		}
		item.CurrentStock = domain.Round3(item.CurrentStock)
		item.IsLowStock = item.LowStockThreshold > 0 && item.CurrentStock <= item.LowStockThreshold
		items = append(items, item)
	}
	return items, rows.Err()
}

func (r *Repository) GetCurrentStockByProduct(ctx context.Context, productID string) (float64, error) {
	var stock float64
	err := r.db.QueryRowContext(ctx, `SELECT COALESCE(SUM(CASE WHEN type='entrada' THEN quantity ELSE -quantity END), 0) FROM movements WHERE product_id = ?`, productID).Scan(&stock)
	if err != nil {
		return 0, err
	}
	return domain.Round3(stock), nil
}

func (r *Repository) GetEntriesByPeriod(ctx context.Context, from time.Time, to time.Time) ([]domain.Movement, error) {
	return r.listByTypeWithPeriod(ctx, domain.MovementTypeEntry, from, to)
}

func (r *Repository) GetEntriesByOrigin(ctx context.Context, from time.Time, to time.Time) ([]domain.GroupedTotal, error) {
	query := `SELECT o.id, o.name, COALESCE(SUM(m.quantity), 0)
		FROM movements m
		JOIN origins o ON o.id = m.origin_id
		WHERE m.type = 'entrada' AND m.movement_date >= ? AND m.movement_date <= ?
		GROUP BY o.id, o.name
		ORDER BY o.name`
	return r.listGrouped(ctx, query, from, to)
}

func (r *Repository) GetExitsByPeriod(ctx context.Context, from time.Time, to time.Time) ([]domain.Movement, error) {
	return r.listByTypeWithPeriod(ctx, domain.MovementTypeExit, from, to)
}

func (r *Repository) GetExitsByDestination(ctx context.Context, from time.Time, to time.Time) ([]domain.GroupedTotal, error) {
	query := `SELECT d.id, d.name, COALESCE(SUM(m.quantity), 0)
		FROM movements m
		JOIN destinations d ON d.id = m.destination_id
		WHERE m.type = 'saida' AND m.movement_date >= ? AND m.movement_date <= ?
		GROUP BY d.id, d.name
		ORDER BY d.name`
	return r.listGrouped(ctx, query, from, to)
}

func (r *Repository) GetMovementHistory(ctx context.Context, from time.Time, to time.Time) ([]domain.Movement, error) {
	return r.ListMovements(ctx, application.ListMovementsFilter{DateFrom: &from, DateTo: &to})
}

func (r *Repository) ListBasketTemplates(ctx context.Context) ([]domain.BasketTemplate, error) {
	rows, err := r.db.QueryContext(ctx, `SELECT id, name, active, created_at, updated_at FROM basket_templates ORDER BY name`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	templates := make([]domain.BasketTemplate, 0)
	templateIDs := make([]string, 0)
	for rows.Next() {
		var item domain.BasketTemplate
		if err := rows.Scan(&item.ID, &item.Name, &item.Active, &item.CreatedAt, &item.UpdatedAt); err != nil {
			return nil, err
		}
		templates = append(templates, item)
		templateIDs = append(templateIDs, item.ID)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	if len(templates) == 0 {
		return templates, nil
	}
	itemsByTemplate, err := r.getBasketTemplateItemsByIDs(ctx, templateIDs)
	if err != nil {
		return nil, err
	}
	for i := range templates {
		templates[i].Items = itemsByTemplate[templates[i].ID]
	}
	return templates, nil
}

func (r *Repository) GetBasketTemplate(ctx context.Context, id string) (*domain.BasketTemplate, error) {
	var item domain.BasketTemplate
	err := r.db.QueryRowContext(ctx, `SELECT id, name, active, created_at, updated_at FROM basket_templates WHERE id = ?`, id).
		Scan(&item.ID, &item.Name, &item.Active, &item.CreatedAt, &item.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil, domain.NewAppError(domain.ErrorNotFound, "basket template not found", nil)
	}
	if err != nil {
		return nil, err
	}

	itemsByTemplate, err := r.getBasketTemplateItemsByIDs(ctx, []string{id})
	if err != nil {
		return nil, err
	}
	item.Items = itemsByTemplate[id]
	return &item, nil
}

func (r *Repository) CreateBasketTemplate(ctx context.Context, item domain.BasketTemplate) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	_, err = tx.ExecContext(ctx, `INSERT INTO basket_templates(id, name, active, created_at, updated_at) VALUES (?, ?, ?, ?, ?)`,
		item.ID, item.Name, item.Active, item.CreatedAt, item.UpdatedAt)
	if err != nil {
		return err
	}
	for _, basketItem := range item.Items {
		_, err := tx.ExecContext(ctx, `INSERT INTO basket_template_items(id, basket_template_id, product_id, quantity_per_basket, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?)`,
			basketItem.ID, item.ID, basketItem.ProductID, basketItem.QuantityPerBasket, basketItem.CreatedAt, basketItem.UpdatedAt)
		if isForeignKeyError(err) {
			return domain.NewAppError(domain.ErrorValidation, "invalid product for basket template", nil)
		}
		if err != nil {
			return err
		}
	}
	return tx.Commit()
}

func (r *Repository) UpdateBasketTemplate(ctx context.Context, item domain.BasketTemplate) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	result, err := tx.ExecContext(ctx, `UPDATE basket_templates SET name = ?, updated_at = ? WHERE id = ?`, item.Name, item.UpdatedAt, item.ID)
	if err != nil {
		return err
	}
	if err := requireRowsAffected(result, "basket template not found"); err != nil {
		return err
	}

	if _, err := tx.ExecContext(ctx, `DELETE FROM basket_template_items WHERE basket_template_id = ?`, item.ID); err != nil {
		return err
	}
	for _, basketItem := range item.Items {
		_, err := tx.ExecContext(ctx, `INSERT INTO basket_template_items(id, basket_template_id, product_id, quantity_per_basket, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?)`,
			basketItem.ID, item.ID, basketItem.ProductID, basketItem.QuantityPerBasket, basketItem.CreatedAt, basketItem.UpdatedAt)
		if isForeignKeyError(err) {
			return domain.NewAppError(domain.ErrorValidation, "invalid product for basket template", nil)
		}
		if err != nil {
			return err
		}
	}
	return tx.Commit()
}

func (r *Repository) DeactivateBasketTemplate(ctx context.Context, id string) error {
	result, err := r.db.ExecContext(ctx, `UPDATE basket_templates SET active = 0, updated_at = ? WHERE id = ?`, time.Now(), id)
	if err != nil {
		return err
	}
	return requireRowsAffected(result, "basket template not found")
}

func (r *Repository) ListFamilies(ctx context.Context) ([]domain.Family, error) {
	rows, err := r.db.QueryContext(ctx, `SELECT id, name, cestas_per_period, active, created_at, updated_at FROM families ORDER BY name`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := make([]domain.Family, 0)
	for rows.Next() {
		var item domain.Family
		if err := rows.Scan(&item.ID, &item.Name, &item.CestasPerPeriod, &item.Active, &item.CreatedAt, &item.UpdatedAt); err != nil {
			return nil, err
		}
		item.CestasPerPeriod = domain.Round3(item.CestasPerPeriod)
		items = append(items, item)
	}
	return items, rows.Err()
}

func (r *Repository) CreateFamily(ctx context.Context, item domain.Family) error {
	_, err := r.db.ExecContext(ctx, `INSERT INTO families(id, name, cestas_per_period, active, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?)`,
		item.ID, item.Name, item.CestasPerPeriod, item.Active, item.CreatedAt, item.UpdatedAt)
	return err
}

func (r *Repository) UpdateFamily(ctx context.Context, item domain.Family) error {
	result, err := r.db.ExecContext(ctx, `UPDATE families SET name = ?, cestas_per_period = ?, updated_at = ? WHERE id = ?`,
		item.Name, item.CestasPerPeriod, item.UpdatedAt, item.ID)
	if err != nil {
		return err
	}
	return requireRowsAffected(result, "family not found")
}

func (r *Repository) DeactivateFamily(ctx context.Context, id string) error {
	result, err := r.db.ExecContext(ctx, `UPDATE families SET active = 0, updated_at = ? WHERE id = ?`, time.Now(), id)
	if err != nil {
		return err
	}
	return requireRowsAffected(result, "family not found")
}

func (r *Repository) GetActiveFamiliesCestas(ctx context.Context) (float64, error) {
	var total float64
	if err := r.db.QueryRowContext(ctx, `SELECT COALESCE(SUM(cestas_per_period), 0) FROM families WHERE active = 1`).Scan(&total); err != nil {
		return 0, err
	}
	return domain.Round3(total), nil
}

func (r *Repository) getBasketTemplateItemsByIDs(ctx context.Context, templateIDs []string) (map[string][]domain.BasketTemplateItem, error) {
	if len(templateIDs) == 0 {
		return map[string][]domain.BasketTemplateItem{}, nil
	}
	placeholders := make([]string, 0, len(templateIDs))
	args := make([]any, 0, len(templateIDs))
	for _, id := range templateIDs {
		placeholders = append(placeholders, "?")
		args = append(args, id)
	}
	query := `SELECT bti.id, bti.basket_template_id, bti.product_id, p.name, p.measure_unit, bti.quantity_per_basket, bti.created_at, bti.updated_at
		FROM basket_template_items bti
		JOIN products p ON p.id = bti.product_id
		WHERE bti.basket_template_id IN (` + strings.Join(placeholders, ",") + `)
		ORDER BY p.name`
	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	itemsByTemplate := map[string][]domain.BasketTemplateItem{}
	for rows.Next() {
		var item domain.BasketTemplateItem
		if err := rows.Scan(&item.ID, &item.BasketTemplateID, &item.ProductID, &item.ProductName, &item.MeasureUnit, &item.QuantityPerBasket, &item.CreatedAt, &item.UpdatedAt); err != nil {
			return nil, err
		}
		item.QuantityPerBasket = domain.Round3(item.QuantityPerBasket)
		itemsByTemplate[item.BasketTemplateID] = append(itemsByTemplate[item.BasketTemplateID], item)
	}
	return itemsByTemplate, rows.Err()
}

func (r *Repository) listByTypeWithPeriod(ctx context.Context, movementType domain.MovementType, from time.Time, to time.Time) ([]domain.Movement, error) {
	query := `SELECT m.id, m.type, m.product_id, p.name, m.quantity, m.movement_date, m.origin_id, o.name, m.destination_id, d.name, m.note, m.created_at
		FROM movements m
		JOIN products p ON p.id = m.product_id
		LEFT JOIN origins o ON o.id = m.origin_id
		LEFT JOIN destinations d ON d.id = m.destination_id
		WHERE m.type = ? AND m.movement_date >= ? AND m.movement_date <= ?
		ORDER BY m.movement_date DESC, m.created_at DESC`

	rows, err := r.db.QueryContext(ctx, query, string(movementType), from, to)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := make([]domain.Movement, 0)
	for rows.Next() {
		var item domain.Movement
		var originID sql.NullString
		var originName sql.NullString
		var destinationID sql.NullString
		var destinationName sql.NullString
		if err := rows.Scan(&item.ID, &item.Type, &item.ProductID, &item.ProductName, &item.Quantity, &item.MovementDate, &originID, &originName, &destinationID, &destinationName, &item.Note, &item.CreatedAt); err != nil {
			return nil, err
		}
		if originID.Valid {
			item.OriginID = &originID.String
		}
		if originName.Valid {
			item.OriginName = &originName.String
		}
		if destinationID.Valid {
			item.DestinationID = &destinationID.String
		}
		if destinationName.Valid {
			item.DestinationName = &destinationName.String
		}
		item.Quantity = domain.Round3(item.Quantity)
		items = append(items, item)
	}
	return items, rows.Err()
}

func (r *Repository) listGrouped(ctx context.Context, query string, from time.Time, to time.Time) ([]domain.GroupedTotal, error) {
	rows, err := r.db.QueryContext(ctx, query, from, to)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := make([]domain.GroupedTotal, 0)
	for rows.Next() {
		var item domain.GroupedTotal
		if err := rows.Scan(&item.ID, &item.Name, &item.Quantity); err != nil {
			return nil, err
		}
		item.Quantity = domain.Round3(item.Quantity)
		items = append(items, item)
	}
	return items, rows.Err()
}

func requireRowsAffected(result sql.Result, message string) error {
	affected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if affected == 0 {
		return domain.NewAppError(domain.ErrorNotFound, message, nil)
	}
	return nil
}

func isForeignKeyError(err error) bool {
	if err == nil {
		return false
	}
	return strings.Contains(strings.ToLower(err.Error()), "foreign key")
}

func insertMovementEntryTx(ctx context.Context, tx *sql.Tx, item domain.Movement) error {
	_, err := tx.ExecContext(
		ctx,
		`INSERT INTO movements(id, type, product_id, quantity, movement_date, origin_id, destination_id, note, created_at)
		 VALUES (?, ?, ?, ?, ?, ?, NULL, ?, ?)`,
		item.ID,
		item.Type,
		item.ProductID,
		item.Quantity,
		item.MovementDate,
		item.OriginID,
		item.Note,
		item.CreatedAt,
	)
	if isForeignKeyError(err) {
		return domain.NewAppError(domain.ErrorValidation, "invalid product or origin", nil)
	}
	return err
}

func insertMovementExitTx(ctx context.Context, tx *sql.Tx, item domain.Movement) error {
	_, err := tx.ExecContext(
		ctx,
		`INSERT INTO movements(id, type, product_id, quantity, movement_date, origin_id, destination_id, note, created_at)
		 VALUES (?, ?, ?, ?, ?, NULL, ?, ?, ?)`,
		item.ID,
		item.Type,
		item.ProductID,
		item.Quantity,
		item.MovementDate,
		item.DestinationID,
		item.Note,
		item.CreatedAt,
	)
	if isForeignKeyError(err) {
		return domain.NewAppError(domain.ErrorValidation, "invalid product or destination", nil)
	}
	return err
}

func (r *Repository) EnsureInterface() application.Repository {
	return r
}

var _ application.Repository = (*Repository)(nil)

var _ = fmt.Sprintf
