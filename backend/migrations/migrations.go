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

// Run executes every *.up.sql migration in lexical order. Files are
// idempotent (CREATE TABLE IF NOT EXISTS / CREATE INDEX IF NOT EXISTS).
//
// We split each file by `;` and run each statement individually against the
// underlying *sql.DB because pgx does not support multi-statement Exec in a
// single prepared call.
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
		if !e.IsDir() && strings.HasSuffix(e.Name(), ".up.sql") {
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

// splitStatements splits a SQL script on `;` boundaries while stripping line
// comments. Migrations are intentionally simple DDL — no PL/pgSQL bodies, no
// string literals containing semicolons.
func splitStatements(sqlText string) []string {
	var out []string
	var cur strings.Builder
	flush := func() {
		s := strings.TrimSpace(cur.String())
		if s != "" {
			out = append(out, s)
		}
		cur.Reset()
	}
	for _, raw := range strings.Split(sqlText, "\n") {
		line := stripLineComment(raw)
		if strings.TrimSpace(line) == "" {
			continue
		}
		cur.WriteString(line)
		cur.WriteByte('\n')
		if strings.HasSuffix(strings.TrimSpace(line), ";") {
			flush()
		}
	}
	flush()
	return out
}

func stripLineComment(s string) string {
	if i := strings.Index(s, "--"); i >= 0 {
		return s[:i]
	}
	return s
}
