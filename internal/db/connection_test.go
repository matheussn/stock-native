package db

import (
	"path/filepath"
	"strings"
	"testing"
)

func TestOpenPathAppliesConnectionPragmas(t *testing.T) {
	dbPath := filepath.Join(t.TempDir(), "pragmas.db")

	conn, err := OpenPath(dbPath)
	if err != nil {
		t.Fatalf("open database: %v", err)
	}
	defer conn.Close()

	var foreignKeys int
	if err := conn.QueryRow("PRAGMA foreign_keys;").Scan(&foreignKeys); err != nil {
		t.Fatalf("read foreign_keys pragma: %v", err)
	}
	if foreignKeys != 1 {
		t.Fatalf("expected foreign_keys=1, got %d", foreignKeys)
	}

	var journalMode string
	if err := conn.QueryRow("PRAGMA journal_mode;").Scan(&journalMode); err != nil {
		t.Fatalf("read journal_mode pragma: %v", err)
	}
	if !strings.EqualFold(journalMode, requiredJournalMode) {
		t.Fatalf("expected journal_mode=%s, got %s", requiredJournalMode, journalMode)
	}

	var busyTimeout int
	if err := conn.QueryRow("PRAGMA busy_timeout;").Scan(&busyTimeout); err != nil {
		t.Fatalf("read busy_timeout pragma: %v", err)
	}
	if busyTimeout != busyTimeoutMillis {
		t.Fatalf("expected busy_timeout=%d, got %d", busyTimeoutMillis, busyTimeout)
	}

	var synchronous int
	if err := conn.QueryRow("PRAGMA synchronous;").Scan(&synchronous); err != nil {
		t.Fatalf("read synchronous pragma: %v", err)
	}
	if synchronous != 1 {
		t.Fatalf("expected synchronous NORMAL (1), got %d", synchronous)
	}
}
