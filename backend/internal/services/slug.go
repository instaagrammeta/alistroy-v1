package services

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"regexp"
	"strconv"
	"strings"
	"unicode"
)

var (
	nonSlug   = regexp.MustCompile(`[^a-z0-9]+`)
	multiDash = regexp.MustCompile(`-+`)
)

// cyrillicMap transliterates Tajik/Russian Cyrillic to Latin.
var cyrillicMap = map[rune]string{
	'а': "a", 'б': "b", 'в': "v", 'г': "g", 'д': "d", 'е': "e", 'ё': "yo",
	'ж': "zh", 'з': "z", 'и': "i", 'й': "y", 'к': "k", 'л': "l", 'м': "m",
	'н': "n", 'о': "o", 'п': "p", 'р': "r", 'с': "s", 'т': "t", 'у': "u",
	'ф': "f", 'х': "kh", 'ц': "ts", 'ч': "ch", 'ш': "sh", 'щ': "shch",
	'ъ': "", 'ы': "y", 'ь': "", 'э': "e", 'ю': "yu", 'я': "ya",
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

// Slugify produces a URL-safe slug, transliterating Cyrillic input.
func Slugify(s string) string {
	s = strings.ToLower(strings.TrimSpace(s))
	s = transliterate(s)
	var b strings.Builder
	for _, r := range s {
		if r <= unicode.MaxASCII {
			b.WriteRune(r)
		}
	}
	out := nonSlug.ReplaceAllString(b.String(), "-")
	out = multiDash.ReplaceAllString(out, "-")
	out = strings.Trim(out, "-")
	if out == "" {
		out = "item"
	}
	if len(out) > 180 {
		out = out[:180]
	}
	return out
}

// uniqueSlug appends -2, -3, ... until exists() returns false.
func uniqueSlug(ctx context.Context, base string, exists func(ctx context.Context, slug string) (bool, error)) (string, error) {
	base = Slugify(base)
	candidate := base
	for i := 2; i < 10000; i++ {
		ok, err := exists(ctx, candidate)
		if err != nil {
			return "", err
		}
		if !ok {
			return candidate, nil
		}
		candidate = base + "-" + strconv.Itoa(i)
	}
	suffix, _ := randomHex(4)
	return base + "-" + suffix, nil
}

func randomHex(n int) (string, error) {
	b := make([]byte, n)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}

func firstNonEmpty(vs ...string) string {
	for _, v := range vs {
		if strings.TrimSpace(v) != "" {
			return v
		}
	}
	return ""
}

func defaultStr(v, def string) string {
	if strings.TrimSpace(v) == "" {
		return def
	}
	return v
}
