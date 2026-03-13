package db

import (
	"database/sql"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type RestoreResult struct {
	RestoredPath       string `json:"restored_path"`
	PreviousBackupPath string `json:"previous_backup_path"`
}

func Backup(conn *sql.DB, activeDBPath, destinationPath string) (string, error) {
	if conn == nil {
		return "", fmt.Errorf("database connection is not initialized")
	}

	activeDBPath, err := normalizeFilePath(activeDBPath)
	if err != nil {
		return "", fmt.Errorf("normalize active database path: %w", err)
	}

	destinationPath, err = normalizeFilePath(destinationPath)
	if err != nil {
		return "", fmt.Errorf("normalize backup destination path: %w", err)
	}

	if samePath(activeDBPath, destinationPath) {
		return "", fmt.Errorf("backup destination must be different from the active database")
	}

	if err := os.MkdirAll(filepath.Dir(destinationPath), 0o755); err != nil {
		return "", fmt.Errorf("create backup directory: %w", err)
	}

	tempPath := destinationPath + ".tmp"
	if err := removeFileIfExists(tempPath); err != nil {
		return "", fmt.Errorf("clear temporary backup file: %w", err)
	}

	if _, err := conn.Exec("VACUUM INTO " + sqliteStringLiteral(tempPath)); err != nil {
		return "", fmt.Errorf("create sqlite backup: %w", err)
	}

	if err := replaceFile(tempPath, destinationPath); err != nil {
		return "", fmt.Errorf("move sqlite backup into place: %w", err)
	}

	return destinationPath, nil
}

func ValidateFile(dbPath string) error {
	dbPath, err := normalizeFilePath(dbPath)
	if err != nil {
		return err
	}

	conn, err := OpenPath(dbPath)
	if err != nil {
		return fmt.Errorf("open database file: %w", err)
	}
	defer conn.Close()

	if err := Migrate(conn); err != nil {
		return fmt.Errorf("migrate database file: %w", err)
	}

	return nil
}

func RestoreFile(activeDBPath, sourceBackupPath string) (RestoreResult, error) {
	activeDBPath, err := normalizeFilePath(activeDBPath)
	if err != nil {
		return RestoreResult{}, fmt.Errorf("normalize active database path: %w", err)
	}

	sourceBackupPath, err = normalizeFilePath(sourceBackupPath)
	if err != nil {
		return RestoreResult{}, fmt.Errorf("normalize backup source path: %w", err)
	}

	if samePath(activeDBPath, sourceBackupPath) {
		return RestoreResult{}, fmt.Errorf("restore source must be different from the active database")
	}

	if _, err := os.Stat(sourceBackupPath); err != nil {
		if os.IsNotExist(err) {
			return RestoreResult{}, fmt.Errorf("backup source file does not exist")
		}
		return RestoreResult{}, fmt.Errorf("stat backup source file: %w", err)
	}

	dir := filepath.Dir(activeDBPath)
	baseName := strings.TrimSuffix(filepath.Base(activeDBPath), filepath.Ext(activeDBPath))
	timestamp := time.Now().Format("20060102-150405")
	tempRestorePath := filepath.Join(dir, baseName+".restore-"+timestamp+".tmp")
	previousBackupPath := filepath.Join(dir, baseName+".before-restore-"+timestamp+".db")

	if err := removeFileIfExists(tempRestorePath); err != nil {
		return RestoreResult{}, fmt.Errorf("clear temporary restore file: %w", err)
	}

	if err := copyFile(sourceBackupPath, tempRestorePath); err != nil {
		return RestoreResult{}, fmt.Errorf("copy restore source to temporary file: %w", err)
	}

	if err := ValidateFile(tempRestorePath); err != nil {
		_ = removeSQLiteSidecars(tempRestorePath)
		_ = removeFileIfExists(tempRestorePath)
		return RestoreResult{}, fmt.Errorf("validate restore source: %w", err)
	}

	if err := removeSQLiteSidecars(tempRestorePath); err != nil {
		_ = removeFileIfExists(tempRestorePath)
		return RestoreResult{}, fmt.Errorf("clear temporary restore sidecars: %w", err)
	}

	hadActiveDatabase := false
	if _, err := os.Stat(activeDBPath); err == nil {
		hadActiveDatabase = true
		if err := removeFileIfExists(previousBackupPath); err != nil {
			_ = removeSQLiteSidecars(tempRestorePath)
			_ = removeFileIfExists(tempRestorePath)
			return RestoreResult{}, fmt.Errorf("clear previous restore backup path: %w", err)
		}
		if err := removeSQLiteSidecars(activeDBPath); err != nil {
			_ = removeSQLiteSidecars(tempRestorePath)
			_ = removeFileIfExists(tempRestorePath)
			return RestoreResult{}, fmt.Errorf("clear active database sidecars before restore: %w", err)
		}
		if err := os.Rename(activeDBPath, previousBackupPath); err != nil {
			_ = removeSQLiteSidecars(tempRestorePath)
			_ = removeFileIfExists(tempRestorePath)
			return RestoreResult{}, fmt.Errorf("move active database to safety backup: %w", err)
		}
	} else if !os.IsNotExist(err) {
		_ = removeSQLiteSidecars(tempRestorePath)
		_ = removeFileIfExists(tempRestorePath)
		return RestoreResult{}, fmt.Errorf("stat active database path: %w", err)
	}

	if err := os.Rename(tempRestorePath, activeDBPath); err != nil {
		if hadActiveDatabase {
			_ = os.Rename(previousBackupPath, activeDBPath)
		}
		_ = removeSQLiteSidecars(tempRestorePath)
		_ = removeFileIfExists(tempRestorePath)
		return RestoreResult{}, fmt.Errorf("move restored database into place: %w", err)
	}

	result := RestoreResult{
		RestoredPath: activeDBPath,
	}
	if hadActiveDatabase {
		result.PreviousBackupPath = previousBackupPath
	}

	return result, nil
}

func normalizeFilePath(filePath string) (string, error) {
	filePath = strings.TrimSpace(filePath)
	if filePath == "" {
		return "", fmt.Errorf("file path is required")
	}

	absolutePath, err := filepath.Abs(filePath)
	if err != nil {
		return "", fmt.Errorf("resolve absolute path: %w", err)
	}

	return filepath.Clean(absolutePath), nil
}

func samePath(left, right string) bool {
	return filepath.Clean(left) == filepath.Clean(right)
}

func replaceFile(sourcePath, destinationPath string) error {
	if err := removeFileIfExists(destinationPath); err != nil {
		return err
	}
	return os.Rename(sourcePath, destinationPath)
}

func removeFileIfExists(filePath string) error {
	err := os.Remove(filePath)
	if err == nil || os.IsNotExist(err) {
		return nil
	}
	return err
}

func copyFile(sourcePath, destinationPath string) error {
	sourceFile, err := os.Open(sourcePath)
	if err != nil {
		return fmt.Errorf("open source file: %w", err)
	}
	defer sourceFile.Close()

	destinationFile, err := os.Create(destinationPath)
	if err != nil {
		return fmt.Errorf("create destination file: %w", err)
	}

	if _, err := io.Copy(destinationFile, sourceFile); err != nil {
		destinationFile.Close()
		return fmt.Errorf("copy file contents: %w", err)
	}

	if err := destinationFile.Sync(); err != nil {
		destinationFile.Close()
		return fmt.Errorf("sync destination file: %w", err)
	}

	if err := destinationFile.Close(); err != nil {
		return fmt.Errorf("close destination file: %w", err)
	}

	return nil
}

func sqliteStringLiteral(value string) string {
	return "'" + strings.ReplaceAll(value, "'", "''") + "'"
}

func removeSQLiteSidecars(dbPath string) error {
	for _, sidecarPath := range sqliteSidecarPaths(dbPath) {
		if err := removeFileIfExists(sidecarPath); err != nil {
			return err
		}
	}
	return nil
}

func sqliteSidecarPaths(dbPath string) []string {
	return []string{
		dbPath + "-wal",
		dbPath + "-shm",
	}
}
