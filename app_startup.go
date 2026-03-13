package main

import (
	"database/sql"
	"fmt"
	"strings"

	"stock/internal/db"
)

type AppStartupStatus struct {
	Ready          bool   `json:"ready"`
	Message        string `json:"message"`
	TechnicalError string `json:"technical_error"`
	DatabasePath   string `json:"database_path"`
}

func (a *App) GetStartupStatus() AppStartupStatus {
	a.mu.RLock()
	defer a.mu.RUnlock()
	return a.startupStatus
}

func (a *App) RetryDatabaseInitialization() AppStartupStatus {
	a.logInfof("database initialization retry requested")
	_ = a.initializeDatabase()
	return a.GetStartupStatus()
}

func (a *App) initializeDatabase() error {
	dbPath, err := a.resolveDatabasePath()
	if err != nil {
		a.markStartupFailure("", err)
		return err
	}

	conn, err := db.OpenPath(dbPath)
	if err != nil {
		a.markStartupFailure(dbPath, err)
		return err
	}

	if err := db.Migrate(conn); err != nil {
		_ = conn.Close()
		a.markStartupFailure(dbPath, err)
		return err
	}

	oldDB := a.setReadyDatabase(conn, dbPath)
	if oldDB != nil && oldDB != conn {
		_ = oldDB.Close()
	}

	a.logInfof("database initialized at: %s", dbPath)

	return nil
}

func (a *App) resolveDatabasePath() (string, error) {
	a.mu.RLock()
	currentPath := strings.TrimSpace(a.dbPath)
	a.mu.RUnlock()

	if currentPath != "" {
		return currentPath, nil
	}

	return db.DBPath()
}

func (a *App) currentDatabase() *sql.DB {
	a.mu.RLock()
	defer a.mu.RUnlock()
	return a.db
}

func (a *App) setReadyDatabase(conn *sql.DB, dbPath string) *sql.DB {
	a.mu.Lock()
	defer a.mu.Unlock()

	oldDB := a.db
	a.db = conn
	a.dbPath = dbPath
	a.startupStatus = AppStartupStatus{
		Ready:        true,
		DatabasePath: dbPath,
	}

	return oldDB
}

func (a *App) markStartupFailure(dbPath string, err error) {
	technical := err.Error()
	userMessage := startupFailureMessage(err)

	a.logErrorf("database initialization failed: %v", err)

	a.mu.Lock()
	currentDB := a.db
	a.db = nil
	a.dbPath = dbPath
	a.startupStatus = AppStartupStatus{
		Ready:          false,
		Message:        userMessage,
		TechnicalError: technical,
		DatabasePath:   dbPath,
	}
	a.mu.Unlock()

	if currentDB != nil {
		_ = currentDB.Close()
	}
}

func (a *App) databaseReadyError() error {
	status := a.GetStartupStatus()
	if status.Ready {
		return nil
	}

	if status.Message != "" {
		return fmt.Errorf("%s", status.Message)
	}

	return fmt.Errorf("database connection is not initialized")
}

func startupFailureMessage(err error) string {
	message := strings.ToLower(err.Error())

	switch {
	case strings.Contains(message, "file is not a database"),
		strings.Contains(message, "database disk image is malformed"),
		strings.Contains(message, "malformed"):
		return "Não foi possível abrir o banco de dados local. O arquivo pode estar corrompido. Tente restaurar um backup."
	case strings.Contains(message, "frozen migration"),
		strings.Contains(message, "migration"),
		strings.Contains(message, "user_version"):
		return "Não foi possível preparar o banco de dados para esta versão do sistema. Revise as migrations ou restaure um backup compatível."
	default:
		return "Não foi possível iniciar o banco de dados local. Verifique o arquivo do banco ou tente restaurar um backup."
	}
}
