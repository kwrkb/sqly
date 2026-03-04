package db

import (
	"strings"
)

// MaskDSN returns a display-safe version of the DSN with passwords masked.
func MaskDSN(dsn string) string {
	if !strings.Contains(dsn, "://") {
		return dsn
	}
	result := dsn
	if idx := strings.Index(result, "://"); idx >= 0 {
		rest := result[idx+3:]
		if atIdx := strings.Index(rest, "@"); atIdx >= 0 {
			userInfo := rest[:atIdx]
			if colonIdx := strings.Index(userInfo, ":"); colonIdx >= 0 {
				result = result[:idx+3] + userInfo[:colonIdx] + ":***" + rest[atIdx:]
			}
		}
	}
	return result
}

// DetectType returns the database type string for a given DSN.
func DetectType(dsn string) string {
	switch {
	case strings.HasPrefix(dsn, "mysql://"):
		return "mysql"
	case strings.HasPrefix(dsn, "postgres://"), strings.HasPrefix(dsn, "postgresql://"):
		return "postgres"
	default:
		return "sqlite"
	}
}

// InitialQuery returns an appropriate initial query for the database type.
func InitialQuery(dbType string) string {
	switch dbType {
	case "mysql":
		return "SELECT VERSION();"
	case "postgres":
		return "SELECT version();"
	default:
		return "SELECT sqlite_version();"
	}
}

// Placeholder returns an appropriate placeholder query for the database type.
func Placeholder(dbType string) string {
	switch dbType {
	case "mysql":
		return "SHOW TABLES;"
	case "postgres":
		return "SELECT tablename FROM pg_tables WHERE schemaname = 'public';"
	default:
		return "SELECT name FROM sqlite_master WHERE type = 'table';"
	}
}

// DisplayName returns a short display name for a DSN.
func DisplayName(dsn string) string {
	switch {
	case strings.HasPrefix(dsn, "mysql://"):
		return extractHost(dsn)
	case strings.HasPrefix(dsn, "postgres://"), strings.HasPrefix(dsn, "postgresql://"):
		return extractHost(dsn)
	default:
		parts := strings.Split(dsn, "/")
		return parts[len(parts)-1]
	}
}

func extractHost(dsn string) string {
	idx := strings.Index(dsn, "://")
	if idx < 0 {
		return dsn
	}
	rest := dsn[idx+3:]
	if atIdx := strings.Index(rest, "@"); atIdx >= 0 {
		rest = rest[atIdx+1:]
	}
	if slashIdx := strings.Index(rest, "/"); slashIdx >= 0 {
		rest = rest[:slashIdx]
	}
	if colonIdx := strings.LastIndex(rest, ":"); colonIdx >= 0 {
		rest = rest[:colonIdx]
	}
	if rest == "" {
		return "localhost"
	}
	return rest
}
