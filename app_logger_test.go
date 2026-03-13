package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestOperationalLoggerWritesStructuredLines(t *testing.T) {
	logPath := filepath.Join(t.TempDir(), "logs", "stock.log")
	logger, err := newOperationalLoggerWithClock(logPath, maxOperationalLogSizeBytes, func() time.Time {
		return time.Date(2026, time.March, 11, 10, 30, 0, 0, time.UTC)
	})
	if err != nil {
		t.Fatalf("create logger: %v", err)
	}
	defer logger.Close()

	logger.Logf("INFO", "movement created: id=%d type=%s", 42, "out")

	content, err := os.ReadFile(logPath)
	if err != nil {
		t.Fatalf("read log file: %v", err)
	}

	text := string(content)
	if !strings.Contains(text, "2026-03-11T10:30:00Z INFO") {
		t.Fatalf("expected timestamp and level in log line, got %q", text)
	}
	if !strings.Contains(text, "movement created: id=42 type=out") {
		t.Fatalf("expected message in log line, got %q", text)
	}
}

func TestOperationalLoggerRotatesWhenFileExceedsLimit(t *testing.T) {
	tempDir := t.TempDir()
	logPath := filepath.Join(tempDir, "stock.log")
	if err := os.WriteFile(logPath, []byte(strings.Repeat("x", 64)), 0o644); err != nil {
		t.Fatalf("seed log file: %v", err)
	}

	logger, err := newOperationalLoggerWithClock(logPath, 16, func() time.Time {
		return time.Date(2026, time.March, 11, 10, 45, 0, 0, time.UTC)
	})
	if err != nil {
		t.Fatalf("create logger: %v", err)
	}
	defer logger.Close()

	rotatedPath := logPath + ".1"
	if _, err := os.Stat(rotatedPath); err != nil {
		t.Fatalf("expected rotated log file: %v", err)
	}

	logger.Logf("INFO", "startup completed")

	currentContent, err := os.ReadFile(logPath)
	if err != nil {
		t.Fatalf("read current log file: %v", err)
	}
	if !strings.Contains(string(currentContent), "startup completed") {
		t.Fatalf("expected new log content after rotation, got %q", string(currentContent))
	}
}
