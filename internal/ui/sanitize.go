package ui

import "strings"

// sanitize strips ANSI escape sequences and control characters from s.
func sanitize(s string) string {
	// Fast path: if no escape characters exist, return as-is
	if !strings.ContainsRune(s, '\x1b') {
		clean := true
		for i := 0; i < len(s); i++ {
			if s[i] < 0x20 && s[i] != '\t' {
				clean = false
				break
			}
		}
		if clean {
			return s
		}
	}

	var b strings.Builder
	b.Grow(len(s))
	i := 0
	for i < len(s) {
		if s[i] == '\x1b' && i+1 < len(s) && s[i+1] == '[' {
			// skip CSI sequence: ESC [ ... final byte
			j := i + 2
			for j < len(s) && s[j] >= 0x20 && s[j] <= 0x3F {
				j++
			}
			if j < len(s) {
				j++ // skip final byte
			}
			i = j
			continue
		}
		if s[i] == '\x1b' {
			// skip other escape sequences (OSC, etc.): ESC ... ST/BEL
			j := i + 1
			for j < len(s) && s[j] != '\x1b' && s[j] != '\a' {
				j++
			}
			if j < len(s) {
				j++
			}
			i = j
			continue
		}
		if s[i] < 0x20 && s[i] != '\t' {
			i++
			continue
		}
		b.WriteByte(s[i])
		i++
	}
	return b.String()
}
