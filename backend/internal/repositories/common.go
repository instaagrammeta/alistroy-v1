package repositories

import (
	"errors"

	"gorm.io/gorm"
)

// IsNotFound is the canonical sentinel check for missing rows.
func IsNotFound(err error) bool {
	return errors.Is(err, gorm.ErrRecordNotFound)
}

// applyPaging clamps page/size and returns the effective offset.
func applyPaging(page, size *int) int {
	if *page < 1 {
		*page = 1
	}
	if *size < 1 || *size > 200 {
		*size = 20
	}
	return (*page - 1) * *size
}
