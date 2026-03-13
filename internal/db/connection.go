package db

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	_ "modernc.org/sqlite"
)

const (
	appDirName          = "stock"
	dbFileName          = "stock.db"
	logDirName          = "logs"
	logFileName         = "stock.log"
	busyTimeoutMillis   = 5000
	synchronousMode     = "NORMAL"
	requiredJournalMode = "WAL"
)

func Open() (*sql.DB, string, error) {
	dbPath, err := DBPath()
	if err != nil {
		return nil, "", err
	}

	conn, err := OpenPath(dbPath)
	if err != nil {
		return nil, "", err
	}

	return conn, dbPath, nil
}

func OpenPath(dbPath string) (*sql.DB, error) {
	dsn := fmt.Sprintf("file:%s", dbPath)
	conn, err := sql.Open("sqlite", dsn)
	if err != nil {
		return nil, fmt.Errorf("open sqlite connection: %w", err)
	}

	if err := applyConnectionPragmas(conn); err != nil {
		conn.Close()
		return nil, err
	}

	if err := conn.Ping(); err != nil {
		conn.Close()
		return nil, fmt.Errorf("ping sqlite connection: %w", err)
	}

	return conn, nil
}

func applyConnectionPragmas(conn *sql.DB) error {
	if _, err := conn.Exec("PRAGMA foreign_keys = ON;"); err != nil {
		return fmt.Errorf("enable foreign keys: %w", err)
	}

	var journalMode string
	if err := conn.QueryRow("PRAGMA journal_mode = WAL;").Scan(&journalMode); err != nil {
		return fmt.Errorf("enable WAL journal mode: %w", err)
	}
	if !strings.EqualFold(journalMode, requiredJournalMode) {
		return fmt.Errorf("unexpected journal mode %q after WAL setup", journalMode)
	}

	if _, err := conn.Exec(fmt.Sprintf("PRAGMA busy_timeout = %d;", busyTimeoutMillis)); err != nil {
		return fmt.Errorf("set busy timeout: %w", err)
	}

	if _, err := conn.Exec("PRAGMA synchronous = NORMAL;"); err != nil {
		return fmt.Errorf("set synchronous mode: %w", err)
	}

	return nil
}

func DBPath() (string, error) {
	appDir, err := AppDataDir()
	if err != nil {
		return "", err
	}

	return filepath.Join(appDir, dbFileName), nil
}

func LogPath() (string, error) {
	appDir, err := AppDataDir()
	if err != nil {
		return "", err
	}

	logDir := filepath.Join(appDir, logDirName)
	if err := os.MkdirAll(logDir, 0o755); err != nil {
		return "", fmt.Errorf("create app log directory: %w", err)
	}

	return filepath.Join(logDir, logFileName), nil
}

func AppDataDir() (string, error) {
	userConfigDir, err := os.UserConfigDir()
	if err != nil {
		return "", fmt.Errorf("resolve user config directory: %w", err)
	}

	appDir := filepath.Join(userConfigDir, appDirName)
	if err := os.MkdirAll(appDir, 0o755); err != nil {
		return "", fmt.Errorf("create app config directory: %w", err)
	}

	return appDir, nil
}
