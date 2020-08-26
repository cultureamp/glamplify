package helper

import (
	"strings"
	"unicode"
)

// ToSnakeCase returns a new string in the format word_word
func ToSnakeCase(s string) string {

	var sb strings.Builder

	in := []rune(strings.TrimSpace(s))
	for i, r := range in {

		if unicode.IsUpper(r) {
			if i > 0 && unicode.IsLower(in[i-1]) {
				sb.WriteRune('_')
			}
			sb.WriteRune(unicode.ToLower(r))

		} else if unicode.IsSpace(r) {
			if !unicode.IsSpace(in[i-1]) {
				sb.WriteRune('_')
			}
		} else {
			sb.WriteRune(r)
		}
	}

	return sb.String()
}

// Redact returns a string with some of the string replaced with *
func Redact(s string) string {
	const stars = 6
	const literals = 4

	l := len(s)
	var b strings.Builder
	b.Grow(l+stars)

	// no matter how long the string, show at least 6 "*"
	for i := 0; i < stars; i++ {
		b.WriteString("*")
	}

	if l <= stars {
		// For small strings, always return "******" don't suffix with any literal chars
		return b.String()
	}

	// For larger strings, redact the first n-chars, and keep the last 4 as is
	r := l - stars // how many of the last chars to print
	if r > literals {
		r = literals // we never print more than the last 4 chars
	}
	r = l - r // index of last chars

	b.WriteString(s[r:])
	return b.String()
}

