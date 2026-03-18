package db

import (
	"path/filepath"
	"strings"
	"testing"
)

func TestValidateMigrationsAcceptsCurrentBaseline(t *testing.T) {
	migrations, err := loadMigrations()
	if err != nil {
		t.Fatalf("load migrations: %v", err)
	}

	if err := validateMigrations(migrations); err != nil {
		t.Fatalf("validate migrations: %v", err)
	}
}

func TestValidateMigrationsRejectsFrozenMigrationChange(t *testing.T) {
	migrations, err := loadMigrations()
	if err != nil {
		t.Fatalf("load migrations: %v", err)
	}

	modified := append([]migration(nil), migrations...)
	modified[0].upSQL = strings.Replace(modified[0].upSQL, "description TEXT,", "description TEXT NOT NULL,", 1)

	err = validateMigrations(modified)
	if err == nil {
		t.Fatal("expected checksum validation error")
	}
	if !strings.Contains(err.Error(), "was modified") {
		t.Fatalf("expected frozen migration error, got: %v", err)
	}
}

func TestValidateMigrationsRejectsVersionGap(t *testing.T) {
	migrations, err := loadMigrations()
	if err != nil {
		t.Fatalf("load migrations: %v", err)
	}

	withGap := append([]migration(nil), migrations...)
	withGap = append(withGap, migration{
		version: len(migrations) + 2,
		name:    "add_indexes",
		upSQL:   "CREATE INDEX idx_product_name ON product(name);",
	})

	err = validateMigrations(withGap)
	if err == nil {
		t.Fatal("expected version sequence validation error")
	}
	if !strings.Contains(err.Error(), "invalid migration version sequence") {
		t.Fatalf("expected version sequence error, got: %v", err)
	}
}

func TestMigrateAppliesLatestEmbeddedSchema(t *testing.T) {
	conn, err := OpenPath(filepath.Join(t.TempDir(), "migration-test.db"))
	if err != nil {
		t.Fatalf("open database: %v", err)
	}
	defer conn.Close()

	if err := Migrate(conn); err != nil {
		t.Fatalf("migrate database: %v", err)
	}

	migrations, err := loadMigrations()
	if err != nil {
		t.Fatalf("load migrations: %v", err)
	}

	version, err := getCurrentVersion(conn)
	if err != nil {
		t.Fatalf("get current version: %v", err)
	}

	if version != len(migrations) {
		t.Fatalf("expected user_version %d, got %d", len(migrations), version)
	}

	expectedIndexes := map[string]struct{}{
		"idx_product_variation_product_active_base":     {},
		"idx_family_assistential_work_active":           {},
		"idx_family_group_assignment_family_started":    {},
		"idx_family_group_assignment_active_family":     {},
		"idx_movement_created":                          {},
		"idx_movement_type_work_created":                {},
		"idx_movement_work_created":                     {},
		"idx_movement_product_item_movement":            {},
		"idx_movement_group_item_movement":              {},
		"idx_movement_group_item_resolution_group_item": {},
	}

	rows, err := conn.Query(`
SELECT name
FROM sqlite_master
WHERE type = 'index'
  AND name LIKE 'idx_%'
`)
	if err != nil {
		t.Fatalf("list indexes: %v", err)
	}
	defer rows.Close()

	foundIndexes := make(map[string]struct{})
	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			t.Fatalf("scan index name: %v", err)
		}
		foundIndexes[name] = struct{}{}
	}
	if err := rows.Err(); err != nil {
		t.Fatalf("iterate indexes: %v", err)
	}

	for indexName := range expectedIndexes {
		if _, ok := foundIndexes[indexName]; !ok {
			t.Fatalf("expected index %q to exist after migrations", indexName)
		}
	}

	var institutionTableCount int
	if err := conn.QueryRow(`
SELECT COUNT(*)
FROM sqlite_master
WHERE type = 'table'
  AND name = 'institution'
`).Scan(&institutionTableCount); err != nil {
		t.Fatalf("check institution table: %v", err)
	}
	if institutionTableCount != 1 {
		t.Fatalf("expected institution table to exist after migrations")
	}

	rows, err = conn.Query(`PRAGMA table_info(movement)`)
	if err != nil {
		t.Fatalf("inspect movement table: %v", err)
	}
	defer rows.Close()

	foundInstitutionColumn := false
	for rows.Next() {
		var cid int
		var name, columnType string
		var notNull, pk int
		var defaultValue any
		if err := rows.Scan(&cid, &name, &columnType, &notNull, &defaultValue, &pk); err != nil {
			t.Fatalf("scan movement column: %v", err)
		}
		if name == "institution_id" {
			foundInstitutionColumn = true
			break
		}
	}
	if err := rows.Err(); err != nil {
		t.Fatalf("iterate movement columns: %v", err)
	}
	if !foundInstitutionColumn {
		t.Fatalf("expected movement.institution_id column to exist after migrations")
	}
}
