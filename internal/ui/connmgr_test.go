package ui

import (
	"testing"

	"github.com/kwrkb/asql/internal/db/sqlite"
)

func TestConnManager(t *testing.T) {
	// Create a temp SQLite adapter for testing
	adapter, err := sqlite.Open(":memory:")
	if err != nil {
		t.Fatalf("failed to open sqlite: %v", err)
	}
	defer adapter.Close()

	cm := newConnManager("test", ":memory:", adapter)

	t.Run("Active returns initial adapter", func(t *testing.T) {
		if cm.Active() != adapter {
			t.Error("expected initial adapter")
		}
	})

	t.Run("ActiveName returns initial name", func(t *testing.T) {
		if cm.ActiveName() != "test" {
			t.Errorf("expected 'test', got %q", cm.ActiveName())
		}
	})

	t.Run("IsConnected for initial DSN", func(t *testing.T) {
		if !cm.IsConnected(":memory:") {
			t.Error("expected initial DSN to be connected")
		}
		if cm.IsConnected("other.db") {
			t.Error("unexpected DSN should not be connected")
		}
	})

	t.Run("IsActive for initial DSN", func(t *testing.T) {
		if !cm.IsActive(":memory:") {
			t.Error("expected initial DSN to be active")
		}
		if cm.IsActive("other.db") {
			t.Error("unexpected DSN should not be active")
		}
	})

	t.Run("Switch to same DSN reuses connection", func(t *testing.T) {
		err := cm.Switch("renamed", ":memory:")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(cm.conns) != 1 {
			t.Errorf("expected 1 connection, got %d", len(cm.conns))
		}
		if cm.ActiveName() != "renamed" {
			t.Errorf("expected name 'renamed', got %q", cm.ActiveName())
		}
	})

	t.Run("Switch to new DSN opens new connection", func(t *testing.T) {
		// Use a file-based temp DB to avoid :memory: conflict
		err := cm.Switch("second", "file::memory:?cache=shared&_fk=1")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(cm.conns) != 2 {
			t.Errorf("expected 2 connections, got %d", len(cm.conns))
		}
		if cm.ActiveName() != "second" {
			t.Errorf("expected name 'second', got %q", cm.ActiveName())
		}
		if cm.IsActive(":memory:") {
			t.Error("first connection should not be active")
		}
	})

	t.Run("Switch back to first connection", func(t *testing.T) {
		err := cm.Switch("test", ":memory:")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if cm.Active() != adapter {
			t.Error("expected to switch back to first adapter")
		}
		if len(cm.conns) != 2 {
			t.Errorf("expected 2 connections (cached), got %d", len(cm.conns))
		}
	})

	t.Run("CloseAll closes all connections", func(t *testing.T) {
		cm.CloseAll()
		if len(cm.conns) != 0 {
			t.Errorf("expected 0 connections after CloseAll, got %d", len(cm.conns))
		}
	})
}
