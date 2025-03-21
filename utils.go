package diplo

import (
	"iter"
	"slices"
	"strings"
	"unicode"
)

func count[I any](it iter.Seq[I]) int {
	count := 0
	for range it {
		count += 1
	}
	return count
}

func hasStringFold(ss []string, s string) bool {
	return slices.ContainsFunc(ss, func(si string) bool {
		return strings.EqualFold(si, s)
	})
}

func simplify(s string) string {
	sb := strings.Builder{}
	sb.Grow(len(s))
	for _, r := range s {
		if unicode.IsLetter(r) && r < 128 {
			sb.WriteRune(r)
		}
	}
	return strings.ToLower(sb.String())
}
