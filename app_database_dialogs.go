package main

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

func (a *App) GetDatabasePath() string {
	dbPath, err := a.resolveDatabasePath()
	if err != nil {
		return ""
	}
	return dbPath
}

func (a *App) GetLogFilePath() string {
	logPath, err := a.resolveLogPath()
	if err != nil {
		return ""
	}
	return logPath
}

func (a *App) ChooseBackupDestination() (string, error) {
	if a.ctx == nil {
		return "", fmt.Errorf("file dialog requires app runtime context")
	}

	return runtime.SaveFileDialog(a.ctx, runtime.SaveDialogOptions{
		Title:                "Salvar backup do banco de dados",
		DefaultDirectory:     preferredDialogDirectory(a.GetDatabasePath()),
		DefaultFilename:      "stock-backup-" + time.Now().Format("20060102-150405") + ".db",
		CanCreateDirectories: true,
		Filters: []runtime.FileFilter{
			{DisplayName: "Banco SQLite (*.db)", Pattern: "*.db"},
			{DisplayName: "Todos os arquivos", Pattern: "*.*"},
		},
	})
}

func (a *App) ChooseRestoreSource() (string, error) {
	if a.ctx == nil {
		return "", fmt.Errorf("file dialog requires app runtime context")
	}

	return runtime.OpenFileDialog(a.ctx, runtime.OpenDialogOptions{
		Title:            "Selecionar arquivo de backup",
		DefaultDirectory: preferredDialogDirectory(a.GetDatabasePath()),
		Filters: []runtime.FileFilter{
			{DisplayName: "Banco SQLite (*.db)", Pattern: "*.db"},
			{DisplayName: "Todos os arquivos", Pattern: "*.*"},
		},
	})
}

func preferredDialogDirectory(dbPath string) string {
	directory := filepath.Dir(dbPath)
	if directoryExists(directory) {
		return directory
	}

	homeDir, err := os.UserHomeDir()
	if err == nil && directoryExists(homeDir) {
		return homeDir
	}

	return ""
}

func directoryExists(path string) bool {
	if path == "" {
		return false
	}

	info, err := os.Stat(path)
	if err != nil {
		return false
	}

	return info.IsDir()
}
