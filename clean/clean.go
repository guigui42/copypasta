package clean

import (
	"regexp"
	"strings"
	"unicode"
)

// Clean applies all cleaning rules to the input text and returns plain text.
func Clean(input string) string {
	s := input
	s = stripHTML(s)
	s = removeOpenInNewTab(s)
	s = removeZeroWidthChars(s)
	s = normalizeBullets(s)
	s = normalizeArrows(s)
	s = stripEmojiPrefixes(s)
	s = trimTrailingWhitespace(s)
	s = collapseBlankLines(s)
	s = strings.TrimSpace(s)
	return s
}

// stripHTML removes any HTML tags that may leak through rich text copy.
var htmlTagRe = regexp.MustCompile(`<[^>]*>`)

func stripHTML(s string) string {
	return htmlTagRe.ReplaceAllString(s, "")
}

// removeOpenInNewTab removes "(opens in a new tab)" annotations.
var openInNewTabRe = regexp.MustCompile(`\s*\(opens in a new tab\)`)

func removeOpenInNewTab(s string) string {
	return openInNewTabRe.ReplaceAllString(s, "")
}

// removeZeroWidthChars strips zero-width spaces, BOM, and other invisible chars.
func removeZeroWidthChars(s string) string {
	return strings.Map(func(r rune) rune {
		switch r {
		case '\u200B', '\u200C', '\u200D', '\uFEFF', '\u00AD', '\u2060':
			return -1
		}
		return r
	}, s)
}

// normalizeBullets replaces various bullet characters with a plain dash.
var bulletRe = regexp.MustCompile(`^(\s*)[•◦○◆▪▸►‣⁃]\s*`)

func normalizeBullets(s string) string {
	lines := strings.Split(s, "\n")
	for i, line := range lines {
		lines[i] = bulletRe.ReplaceAllString(line, "${1}- ")
	}
	return strings.Join(lines, "\n")
}

// normalizeArrows replaces → with ->.
func normalizeArrows(s string) string {
	return strings.ReplaceAll(s, "→", "->")
}

// stripEmojiPrefixes removes leading emoji from lines (e.g., "✅ Text" → "Text").
func stripEmojiPrefixes(s string) string {
	lines := strings.Split(s, "\n")
	for i, line := range lines {
		lines[i] = stripLeadingEmoji(line)
	}
	return strings.Join(lines, "\n")
}

func stripLeadingEmoji(line string) string {
	trimmed := strings.TrimLeftFunc(line, unicode.IsSpace)
	if trimmed == "" {
		return line
	}
	leadingSpaces := line[:len(line)-len(trimmed)]

	runes := []rune(trimmed)
	idx := 0
	stripped := false
	for idx < len(runes) {
		r := runes[idx]
		if isEmojiLike(r) {
			idx++
			stripped = true
			continue
		}
		// Skip variation selectors and ZWJ after emoji
		if r == '\uFE0F' || r == '\uFE0E' || r == '\u200D' {
			idx++
			continue
		}
		break
	}
	if stripped {
		rest := strings.TrimLeftFunc(string(runes[idx:]), unicode.IsSpace)
		if rest == "" {
			return leadingSpaces
		}
		return leadingSpaces + rest
	}
	return line
}

func isEmojiLike(r rune) bool {
	// Common emoji ranges
	if r >= 0x1F300 && r <= 0x1FAF8 { // Misc Symbols, Emoticons, etc.
		return true
	}
	if r >= 0x2600 && r <= 0x27BF { // Misc symbols, Dingbats
		return true
	}
	if r >= 0x2700 && r <= 0x27BF { // Dingbats
		return true
	}
	if r >= 0xFE00 && r <= 0xFE0F { // Variation selectors
		return true
	}
	if r >= 0x1F900 && r <= 0x1F9FF { // Supplemental Symbols
		return true
	}
	// Check mark, ballot box, etc.
	switch r {
	case '✅', '❌', '⚠', '📷', '🔗', '💡', '🎯', '🚀', '✓', '✔', '☑', '☐', '📌', '📝', '🔒', '🔑':
		return true
	}
	return false
}

// trimTrailingWhitespace removes trailing spaces/tabs from each line.
func trimTrailingWhitespace(s string) string {
	lines := strings.Split(s, "\n")
	for i, line := range lines {
		lines[i] = strings.TrimRightFunc(line, unicode.IsSpace)
	}
	return strings.Join(lines, "\n")
}

// collapseBlankLines reduces 3+ consecutive blank lines to 1.
func collapseBlankLines(s string) string {
	lines := strings.Split(s, "\n")
	var result []string
	blankCount := 0
	for _, line := range lines {
		if strings.TrimSpace(line) == "" {
			blankCount++
			if blankCount <= 1 {
				result = append(result, "")
			}
		} else {
			blankCount = 0
			result = append(result, line)
		}
	}
	return strings.Join(result, "\n")
}
