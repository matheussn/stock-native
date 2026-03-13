package main

import (
	"os"
	"path/filepath"
	"testing"

	"stock/internal/db"
)

func TestAppBackupAndRestoreDatabase(t *testing.T) {
	tempDir := t.TempDir()
	activeDBPath := filepath.Join(tempDir, "active.db")

	activeDB, err := db.OpenPath(activeDBPath)
	if err != nil {
		t.Fatalf("open active database: %v", err)
	}

	if err := db.Migrate(activeDB); err != nil {
		activeDB.Close()
		t.Fatalf("migrate active database: %v", err)
	}

	app := &App{
		db:     activeDB,
		dbPath: activeDBPath,
		startupStatus: AppStartupStatus{
			Ready:        true,
			DatabasePath: activeDBPath,
		},
	}
	defer func() {
		if app.db != nil {
			_ = app.db.Close()
		}
	}()

	if _, err := app.CreateAssistentialWork("Current Work", ""); err != nil {
		t.Fatalf("seed active database: %v", err)
	}

	backupPath := filepath.Join(tempDir, "backups", "stock-backup.db")
	backedUpPath, err := app.BackupDatabase(backupPath)
	if err != nil {
		t.Fatalf("backup database: %v", err)
	}

	if _, err := os.Stat(backedUpPath); err != nil {
		t.Fatalf("stat backup database: %v", err)
	}

	backupDB, err := db.OpenPath(backedUpPath)
	if err != nil {
		t.Fatalf("open backup database: %v", err)
	}
	defer backupDB.Close()

	var backedUpCount int
	if err := backupDB.QueryRow(`SELECT COUNT(*) FROM assistential_work`).Scan(&backedUpCount); err != nil {
		t.Fatalf("count assistential works in backup: %v", err)
	}
	if backedUpCount != 1 {
		t.Fatalf("expected 1 assistential work in backup, got %d", backedUpCount)
	}

	restoreSourcePath := filepath.Join(tempDir, "restore-source.db")
	restoreSourceDB, err := db.OpenPath(restoreSourcePath)
	if err != nil {
		t.Fatalf("open restore source database: %v", err)
	}

	if err := db.Migrate(restoreSourceDB); err != nil {
		restoreSourceDB.Close()
		t.Fatalf("migrate restore source database: %v", err)
	}

	if _, err := restoreSourceDB.Exec(`INSERT INTO assistential_work (name, description, is_active) VALUES ('Restored Work', '', 1)`); err != nil {
		restoreSourceDB.Close()
		t.Fatalf("seed restore source database: %v", err)
	}

	if err := restoreSourceDB.Close(); err != nil {
		t.Fatalf("close restore source database: %v", err)
	}

	restoreResult, err := app.RestoreDatabase(restoreSourcePath)
	if err != nil {
		t.Fatalf("restore database: %v", err)
	}

	if restoreResult.PreviousBackupPath == "" {
		t.Fatal("expected previous backup path after restore")
	}

	if _, err := os.Stat(restoreResult.PreviousBackupPath); err != nil {
		t.Fatalf("stat previous backup path: %v", err)
	}

	if app.db == nil {
		t.Fatal("expected app database connection to be reopened after restore")
	}

	works, err := app.ListAssistentialWorks(false)
	if err != nil {
		t.Fatalf("list assistential works after restore: %v", err)
	}

	if len(works) != 1 {
		t.Fatalf("expected 1 assistential work after restore, got %d", len(works))
	}

	if works[0].Name != "Restored Work" {
		t.Fatalf("expected restored assistential work, got %q", works[0].Name)
	}
}
