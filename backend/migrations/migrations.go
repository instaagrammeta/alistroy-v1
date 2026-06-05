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
func Run(db *gorm.DB) error {
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
		sqlText := string(b)
		if strings.TrimSpace(sqlText) == "" {
			continue
		}
		if err := db.Exec(sqlText).Error; err != nil {
			return fmt.Errorf("apply %s: %w", name, err)
		}
	}
	return nil
}
