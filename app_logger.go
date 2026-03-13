package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"stock/internal/db"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

const maxOperationalLogSizeBytes int64 = 5 * 1024 * 1024

type operationalLogger struct {
	mu    sync.Mutex
	file  *os.File
	path  string
	clock func() time.Time
}

func newOperationalLogger(path string) (*operationalLogger, error) {
	return newOperationalLoggerWithClock(path, maxOperationalLogSizeBytes, time.Now)
}

func newOperationalLoggerWithClock(path string, maxSizeBytes int64, clock func() time.Time) (*operationalLogger, error) {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return nil, fmt.Errorf("create operational log directory: %w", err)
	}

	if err := rotateOperationalLogFileIfNeeded(path, maxSizeBytes); err != nil {
		return nil, err
	}

	file, err := os.OpenFile(path, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0o644)
	if err != nil {
		return nil, fmt.Errorf("open operational log file: %w", err)
	}

	return &operationalLogger{
		file:  file,
		path:  path,
		clock: clock,
	}, nil
}

func rotateOperationalLogFileIfNeeded(path string, maxSizeBytes int64) error {
	info, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return fmt.Errorf("stat operational log file: %w", err)
	}

	if info.Size() < maxSizeBytes {
		return nil
	}

	backupPath := path + ".1"
	if err := os.Remove(backupPath); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("remove previous operational log backup: %w", err)
	}

	if err := os.Rename(path, backupPath); err != nil {
		return fmt.Errorf("rotate operational log file: %w", err)
	}

	return nil
}

func (l *operationalLogger) Logf(level, format string, args ...any) {
	if l == nil {
		return
	}

	message := fmt.Sprintf(format, args...)
	line := fmt.Sprintf("%s %-5s %s\n", l.clock().Format(time.RFC3339), level, message)

	l.mu.Lock()
	defer l.mu.Unlock()

	if l.file == nil {
		return
	}

	_, _ = l.file.WriteString(line)
}

func (l *operationalLogger) Close() error {
	if l == nil {
		return nil
	}

	l.mu.Lock()
	defer l.mu.Unlock()

	if l.file == nil {
		return nil
	}

	err := l.file.Close()
	l.file = nil
	return err
}

func (a *App) initializeOperationalLogger() error {
	logPath, err := a.resolveLogPath()
	if err != nil {
		return err
	}

	logger, err := newOperationalLogger(logPath)
	if err != nil {
		return err
	}

	a.mu.Lock()
	oldLogger := a.logger
	a.logger = logger
	a.logPath = logPath
	a.mu.Unlock()

	if oldLogger != nil && oldLogger != logger {
		_ = oldLogger.Close()
	}

	a.logInfof("operational logger initialized at: %s", logPath)
	return nil
}

func (a *App) resolveLogPath() (string, error) {
	a.mu.RLock()
	currentPath := strings.TrimSpace(a.logPath)
	a.mu.RUnlock()

	if currentPath != "" {
		return currentPath, nil
	}

	return db.LogPath()
}

func (a *App) currentLogger() *operationalLogger {
	a.mu.RLock()
	defer a.mu.RUnlock()
	return a.logger
}

func (a *App) closeOperationalLogger() {
	a.mu.Lock()
	logger := a.logger
	a.logger = nil
	a.mu.Unlock()

	if logger != nil {
		if err := logger.Close(); err != nil && a.ctx != nil {
			runtime.LogErrorf(a.ctx, "failed to close operational log file: %v", err)
		}
	}
}

func (a *App) logInfof(format string, args ...any) {
	a.logf("INFO", format, args...)
	if a.ctx != nil {
		runtime.LogInfof(a.ctx, format, args...)
	}
}

func (a *App) logErrorf(format string, args ...any) {
	a.logf("ERROR", format, args...)
	if a.ctx != nil {
		runtime.LogErrorf(a.ctx, format, args...)
	}
}

func (a *App) logf(level, format string, args ...any) {
	logger := a.currentLogger()
	if logger != nil {
		logger.Logf(level, format, args...)
	}
}

func (a *App) logMutationSuccess(action, details string, args ...any) {
	if strings.TrimSpace(details) == "" {
		a.logInfof("%s", action)
		return
	}

	a.logInfof("%s: %s", action, fmt.Sprintf(details, args...))
}

func (a *App) logMutationFailure(action string, err error, details string, args ...any) {
	if strings.TrimSpace(details) == "" {
		a.logErrorf("%s failed: %v", action, err)
		return
	}

	a.logErrorf("%s failed: %s; err=%v", action, fmt.Sprintf(details, args...), err)
}
