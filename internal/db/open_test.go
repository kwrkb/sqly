package db

import "testing"

func TestDetectType(t *testing.T) {
	tests := []struct {
		dsn  string
		want string
	}{
		{"mysql://user:pass@host/db", "mysql"},
		{"postgres://user:pass@host/db", "postgres"},
		{"postgresql://user:pass@host/db", "postgres"},
		{"test.db", "sqlite"},
		{"/path/to/file.db", "sqlite"},
	}
	for _, tt := range tests {
		if got := DetectType(tt.dsn); got != tt.want {
			t.Errorf("DetectType(%q) = %q, want %q", tt.dsn, got, tt.want)
		}
	}
}

func TestDisplayName(t *testing.T) {
	tests := []struct {
		dsn  string
		want string
	}{
		{"test.db", "test.db"},
		{"/path/to/data.db", "data.db"},
		{"mysql://user:pass@myhost:3306/mydb", "myhost"},
		{"postgres://user:pass@pghost:5432/pgdb", "pghost"},
		{"postgres://user@localhost/db", "localhost"},
	}
	for _, tt := range tests {
		if got := DisplayName(tt.dsn); got != tt.want {
			t.Errorf("DisplayName(%q) = %q, want %q", tt.dsn, got, tt.want)
		}
	}
}

func TestMaskDSN(t *testing.T) {
	tests := []struct {
		dsn  string
		want string
	}{
		{"test.db", "test.db"},
		{"mysql://user:secret@host/db", "mysql://user:***@host/db"},
		{"postgres://admin:p4ss@host:5432/db", "postgres://admin:***@host:5432/db"},
		{"sqlite.db", "sqlite.db"},
	}
	for _, tt := range tests {
		if got := MaskDSN(tt.dsn); got != tt.want {
			t.Errorf("MaskDSN(%q) = %q, want %q", tt.dsn, got, tt.want)
		}
	}
}
