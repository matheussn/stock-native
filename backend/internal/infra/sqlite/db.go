package sqlite

import (
	"context"
	"database/sql"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	_ "modernc.org/sqlite"

	"stock/backend/internal/infra/migrations"
)

type DBManager struct {
	db     *sql.DB
	dbPath string
}

func NewDBManager(dbPath string) (*DBManager, error) {
	if err := os.MkdirAll(filepath.Dir(dbPath), 0o755); err != nil {
		return nil, err
	}
	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return nil, err
	}

	pragmas := []string{
		"PRAGMA foreign_keys = ON;",
		"PRAGMA journal_mode = WAL;",
		"PRAGMA busy_timeout = 5000;",
	}
	for _, pragma := range pragmas {
		if _, execErr := db.Exec(pragma); execErr != nil {
			_ = db.Close()
			return nil, execErr
		}
	}

	if err := runMigrations(db); err != nil {
		_ = db.Close()
		return nil, err
	}

	return &DBManager{db: db, dbPath: dbPath}, nil
}

func (m *DBManager) DB() *sql.DB {
	return m.db
}

func (m *DBManager) Close() error {
	if m.db == nil {
		return nil
	}
	return m.db.Close()
}

func (m *DBManager) ExportBackup(ctx context.Context, destination string) error {
	if strings.TrimSpace(destination) == "" {
		return fmt.Errorf("destination is required")
	}

	if _, err := m.db.ExecContext(ctx, "PRAGMA wal_checkpoint(FULL);"); err != nil {
		return err
	}
	if _, err := m.db.ExecContext(ctx, fmt.Sprintf("VACUUM INTO '%s';", escapeSQLitePath(destination))); err != nil {
		return err
	}
	return nil
}

func (m *DBManager) ImportBackup(ctx context.Context, source string) error {
	if strings.TrimSpace(source) == "" {
		return fmt.Errorf("source is required")
	}
	if _, err := os.Stat(source); err != nil {
		return err
	}

	backupName := fmt.Sprintf("app_before_import_%s.db", time.Now().Format("20060102_150405"))
	backupPath := filepath.Join(filepath.Dir(m.dbPath), backupName)

	if err := m.ExportBackup(ctx, backupPath); err != nil {
		return err
	}

	if err := m.db.Close(); err != nil {
		return err
	}

	data, err := os.ReadFile(source)
	if err != nil {
		return err
	}
	if err := os.WriteFile(m.dbPath, data, 0o644); err != nil {
		return err
	}

	newDB, err := sql.Open("sqlite", m.dbPath)
	if err != nil {
		return err
	}
	m.db = newDB

	pragmas := []string{
		"PRAGMA foreign_keys = ON;",
		"PRAGMA journal_mode = WAL;",
		"PRAGMA busy_timeout = 5000;",
	}
	for _, pragma := range pragmas {
		if _, execErr := m.db.ExecContext(ctx, pragma); execErr != nil {
			return execErr
		}
	}
	return runMigrations(m.db)
}

func runMigrations(db *sql.DB) error {
	dirEntries, err := fs.ReadDir(migrations.FS, ".")
	if err != nil {
		return err
	}
	files := make([]string, 0, len(dirEntries))
	for _, entry := range dirEntries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".sql") {
			continue
		}
		files = append(files, entry.Name())
	}
	sort.Strings(files)

	var currentVersion int
	if err := db.QueryRow("PRAGMA user_version;").Scan(&currentVersion); err != nil {
		return err
	}

	for i, filename := range files {
		version := i + 1
		if version <= currentVersion {
			continue
		}

		content, readErr := migrations.FS.ReadFile(filename)
		if readErr != nil {
			return readErr
		}
		if _, execErr := db.Exec(string(content)); execErr != nil {
			return fmt.Errorf("migration %s failed: %w", filename, execErr)
		}
		if _, execErr := db.Exec(fmt.Sprintf("PRAGMA user_version = %d;", version)); execErr != nil {
			return execErr
		}
	}

	return nil
}

func escapeSQLitePath(path string) string {
	return strings.ReplaceAll(path, "'", "''")
}
