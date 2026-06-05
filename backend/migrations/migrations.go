// Package migrations embeds the SQL migration files and exposes a runner.
package migrations

import (
	"embed"
	"fmt"
	"sort"
	"strings"

	"gorm.io/gorm"
)

//go:embed *.sql
var FS embed.FS

// Run executes every *.up.sql migration in lexical order.
// Migrations are written to be idempotent (CREATE TABLE IF NOT EXISTS, etc.).
//
// The pgx driver does not support multiple SQL statements in a single
// prepared-statement Exec, so we split each migration file into individual
// statements (one per terminating semicolon) and execute them one by one
// against the underlying *sql.DB (bypassing GORM's prepared-statement cache).
func Run(db *gorm.DB) error {
	sqlDB, err := db.DB()
	if err != nil {
		return fmt.Errorf("get sql.DB: %w", err)
	}

	entries, err := FS.ReadDir(".")
	if err != nil {
		return fmt.Errorf("read migrations: %w", err)
	}
	var ups []string
	for _, e := range entries {
		if e.IsDir() {
			continue
		}
		if strings.HasSuffix(e.Name(), ".up.sql") {
			ups = append(ups, e.Name())
		}
	}
	sort.Strings(ups)

	for _, name := range ups {
		b, err := FS.ReadFile(name)
		if err != nil {
			return fmt.Errorf("read %s: %w", name, err)
		}
		stmts := splitStatements(string(b))
		for i, stmt := range stmts {
			if _, err := sqlDB.Exec(stmt); err != nil {
				return fmt.Errorf("apply %s [stmt %d]: %w", name, i+1, err)
			}
		}
	}
	return nil
}

// splitStatements splits a SQL script into individual statements separated by
// `;`. Line comments (-- ...) are stripped. This is intentionally simple: our
// migrations contain plain DDL only — no string literals containing `;`,
// no PL/pgSQL function bodies, etc.
func splitStatements(sqlText string) []string {
	var stmts []string
	var current strings.Builder

	flush := func() {
		s := strings.TrimSpace(current.String())
		if s != "" {
			stmts = append(stmts, s)
		}
		current.Reset()
	}

	for _, rawLine := range strings.Split(sqlText, "\n") {
		line := stripLineComment(rawLine)
		trimmed := strings.TrimSpace(line)
		if trimmed == "" {
			continue
		}
		current.WriteString(line)
		current.WriteByte('\n')
		if strings.HasSuffix(trimmed, ";") {
			flush()
		}
	}
	flush()
	return stmts
}

// stripLineComment removes a trailing `-- ...` portion from a SQL line.
func stripLineComment(s string) string {
	if i := strings.Index(s, "--"); i >= 0 {
		return s[:i]
	}
	return s
}
