package services

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"regexp"
	"strings"
	"unicode"
)

// Common service-layer errors. Handlers map these to HTTP responses.
var (
	ErrNotFound           = errors.New("not found")
	ErrConflict           = errors.New("conflict")
	ErrValidation         = errors.New("validation failed")
	ErrUnauthorized       = errors.New("unauthorized")
	ErrForbidden          = errors.New("forbidden")
	ErrInvalidCredentials = errors.New("invalid credentials")
)

// Slugify produces a URL-safe slug from arbitrary input. It transliterates
// Tajik/Russian Cyrillic to Latin and drops everything that isn't [a-z0-9].
func Slugify(s string) string {
	s = strings.ToLower(strings.TrimSpace(s))
	s = transliterate(s)

	var b strings.Builder
	for _, r := range s {
		if r > unicode.MaxASCII {
			continue
		}
		b.WriteRune(r)
	}
	out := b.String()
	out = nonSlug.ReplaceAllString(out, "-")
	out = multiDash.ReplaceAllString(out, "-")
	out = strings.Trim(out, "-")
	if out == "" {
		out = "item"
	}
	if len(out) > 200 {
		out = out[:200]
	}
	return out
}

var (
	nonSlug   = regexp.MustCompile(`[^a-z0-9]+`)
	multiDash = regexp.MustCompile(`-+`)
)

var cyrillicMap = map[rune]string{
	'а': "a", 'б': "b", 'в': "v", 'г': "g", 'д': "d", 'е': "e", 'ё': "yo",
	'ж': "zh", 'з': "z", 'и': "i", 'й': "y", 'к': "k", 'л': "l", 'м': "m",
	'н': "n", 'о': "o", 'п': "p", 'р': "r", 'с': "s", 'т': "t", 'у': "u",
	'ф': "f", 'х': "kh", 'ц': "ts", 'ч': "ch", 'ш': "sh", 'щ': "shch",
	'ъ': "", 'ы': "y", 'ь': "", 'э': "e", 'ю': "yu", 'я': "ya",
	// Tajik specific
	'ӣ': "i", 'ӯ': "u", 'қ': "q", 'ҳ': "h", 'ҷ': "j", 'ғ': "gh",
}

func transliterate(s string) string {
	var b strings.Builder
	b.Grow(len(s))
	for _, r := range s {
		if v, ok := cyrillicMap[unicode.ToLower(r)]; ok {
			b.WriteString(v)
		} else {
			b.WriteRune(r)
		}
	}
	return b.String()
}

// RandomToken returns a hex-encoded random token of the given byte length.
func RandomToken(n int) (string, error) {
	b := make([]byte, n)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}
