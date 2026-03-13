package db

import (
	"crypto/sha256"
	"database/sql"
	"embed"
	"encoding/hex"
	"fmt"
	"sort"
	"strconv"
	"strings"
)

type migration struct {
	version  int
	name     string
	upSQL    string
	frozen   bool
	checksum string
}

const initialSchemaChecksum = "5ed8d8b1c65bb8b3eb376b6d2d513ecf551e117e8ed851609ecac45d73eed1dd"

var frozenMigrationChecksums = map[int]string{
	1: initialSchemaChecksum,
}

//go:embed migrations
var migrationFiles embed.FS

func Migrate(conn *sql.DB) error {
	migrations, err := loadMigrations()
	if err != nil {
		return err
	}

	if err := validateMigrations(migrations); err != nil {
		return err
	}

	currentVersion, err := getCurrentVersion(conn)
	if err != nil {
		return err
	}

	for _, m := range migrations {
		if m.version <= currentVersion {
			continue
		}

		if err := applyMigration(conn, m); err != nil {
			return fmt.Errorf("apply migration %d (%s): %w", m.version, m.name, err)
		}
	}

	return nil
}

func loadMigrations() ([]migration, error) {
	entries, err := migrationFiles.ReadDir("migrations")
	if err != nil {
		return nil, fmt.Errorf("read embedded migrations: %w", err)
	}

	list := make([]migration, 0, len(entries))
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		if !strings.HasSuffix(entry.Name(), ".sql") {
			continue
		}

		version, name, err := parseMigrationFileName(entry.Name())
		if err != nil {
			return nil, err
		}

		upSQLBytes, err := migrationFiles.ReadFile("migrations/" + entry.Name())
		if err != nil {
			return nil, fmt.Errorf("read embedded migration %q: %w", entry.Name(), err)
		}

		m := migration{
			version: version,
			name:    name,
			upSQL:   string(upSQLBytes),
		}
		if checksum, ok := frozenMigrationChecksums[version]; ok {
			m.frozen = true
			m.checksum = checksum
		}

		list = append(list, m)
	}

	sort.Slice(list, func(i, j int) bool {
		return list[i].version < list[j].version
	})

	if len(list) == 0 {
		return nil, fmt.Errorf("no .sql migration files found in embedded migrations directory")
	}

	return list, nil
}

func parseMigrationFileName(fileName string) (int, string, error) {
	if !strings.HasSuffix(fileName, ".sql") {
		return 0, "", fmt.Errorf("invalid migration file %q: expected .sql extension", fileName)
	}

	baseName := strings.TrimSuffix(fileName, ".sql")
	parts := strings.SplitN(baseName, "_", 2)
	if len(parts) != 2 {
		return 0, "", fmt.Errorf("invalid migration file %q: expected NNN_name.sql", fileName)
	}

	version, err := strconv.Atoi(parts[0])
	if err != nil || version <= 0 {
		return 0, "", fmt.Errorf("invalid migration file %q: invalid version prefix", fileName)
	}

	name := strings.TrimSpace(parts[1])
	if name == "" {
		return 0, "", fmt.Errorf("invalid migration file %q: empty migration name", fileName)
	}

	return version, name, nil
}

func validateMigrations(list []migration) error {
	expectedVersion := 1

	for _, m := range list {
		if m.version != expectedVersion {
			return fmt.Errorf("invalid migration version sequence: expected %d, got %d (%s)", expectedVersion, m.version, m.name)
		}
		if strings.TrimSpace(m.name) == "" {
			return fmt.Errorf("migration %d has empty name", m.version)
		}
		if strings.TrimSpace(m.upSQL) == "" {
			return fmt.Errorf("migration %d (%s) has empty SQL", m.version, m.name)
		}
		if m.frozen {
			if m.checksum == "" {
				return fmt.Errorf("frozen migration %d (%s) is missing checksum", m.version, m.name)
			}

			actualChecksum := migrationChecksum(m.upSQL)
			if actualChecksum != m.checksum {
				return fmt.Errorf(
					"frozen migration %d (%s) was modified; expected checksum %s, got %s",
					m.version,
					m.name,
					m.checksum,
					actualChecksum,
				)
			}
		}

		expectedVersion++
	}

	return nil
}

func getCurrentVersion(conn *sql.DB) (int, error) {
	var version int
	if err := conn.QueryRow("PRAGMA user_version;").Scan(&version); err != nil {
		return 0, fmt.Errorf("read current migration version: %w", err)
	}
	return version, nil
}

func applyMigration(conn *sql.DB, m migration) error {
	tx, err := conn.Begin()
	if err != nil {
		return fmt.Errorf("begin migration transaction: %w", err)
	}

	defer func() {
		_ = tx.Rollback()
	}()

	if _, err := tx.Exec(m.upSQL); err != nil {
		return fmt.Errorf("execute migration SQL: %w", err)
	}

	if _, err := tx.Exec(fmt.Sprintf("PRAGMA user_version = %d;", m.version)); err != nil {
		return fmt.Errorf("update migration version: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit migration transaction: %w", err)
	}

	return nil
}

func migrationChecksum(sql string) string {
	normalized := strings.TrimSpace(strings.ReplaceAll(sql, "\r\n", "\n"))
	sum := sha256.Sum256([]byte(normalized))
	return hex.EncodeToString(sum[:])
}
