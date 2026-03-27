package clean

import (
	"testing"
)

func TestClean_SampleText(t *testing.T) {
	input := `1. How to visualize premium request consumption

A. What individual users can see (self‑service)

Each user can only see their own usage.

Where:

• In the IDE (VS Code, JetBrains, Visual Studio, Xcode, etc.)
   ◦ Copilot icon → quota / usage
• On GitHub.com:
   ◦ https://github.com/settings/billing (opens in a new tab)
   ◦ Scroll to Metered usage
   ◦ Filter by Copilot

This shows:

• How many premium requests they've used this month
• When the counter resets (always the 1st of the month, UTC)

Source: Monitoring your GitHub Copilot usage and entitlements (opens in a new tab)`

	expected := `1. How to visualize premium request consumption

A. What individual users can see (self‑service)

Each user can only see their own usage.

Where:

- In the IDE (VS Code, JetBrains, Visual Studio, Xcode, etc.)
   - Copilot icon -> quota / usage
- On GitHub.com:
   - https://github.com/settings/billing
   - Scroll to Metered usage
   - Filter by Copilot

This shows:

- How many premium requests they've used this month
- When the counter resets (always the 1st of the month, UTC)

Source: Monitoring your GitHub Copilot usage and entitlements`

	got := Clean(input)
	if got != expected {
		t.Errorf("Clean() mismatch.\n\nGOT:\n%s\n\nEXPECTED:\n%s", got, expected)
	}
}

func TestRemoveOpenInNewTab(t *testing.T) {
	cases := []struct {
		in, want string
	}{
		{"click here (opens in a new tab)", "click here"},
		{"link (opens in a new tab) more text", "link more text"},
		{"no annotation here", "no annotation here"},
	}
	for _, c := range cases {
		got := removeOpenInNewTab(c.in)
		if got != c.want {
			t.Errorf("removeOpenInNewTab(%q) = %q, want %q", c.in, got, c.want)
		}
	}
}

func TestNormalizeBullets(t *testing.T) {
	input := "• First item\n   ◦ Sub item\n○ Another"
	expected := "- First item\n   - Sub item\n- Another"
	got := normalizeBullets(input)
	if got != expected {
		t.Errorf("normalizeBullets():\ngot:  %q\nwant: %q", got, expected)
	}
}

func TestStripEmojiPrefixes(t *testing.T) {
	cases := []struct {
		in, want string
	}{
		{"✅ Enterprise owner", "Enterprise owner"},
		{"📷 Screenshot here", "Screenshot here"},
		{"Normal text", "Normal text"},
		{"  ✅ Indented emoji", "  Indented emoji"},
	}
	for _, c := range cases {
		got := stripLeadingEmoji(c.in)
		if got != c.want {
			t.Errorf("stripLeadingEmoji(%q) = %q, want %q", c.in, got, c.want)
		}
	}
}

func TestCollapseBlankLines(t *testing.T) {
	input := "line1\n\n\n\n\nline2\n\nline3"
	expected := "line1\n\nline2\n\nline3"
	got := collapseBlankLines(input)
	if got != expected {
		t.Errorf("collapseBlankLines():\ngot:  %q\nwant: %q", got, expected)
	}
}

func TestStripHTML(t *testing.T) {
	input := "<b>Bold</b> and <a href='x'>link</a>"
	expected := "Bold and link"
	got := stripHTML(input)
	if got != expected {
		t.Errorf("stripHTML(%q) = %q, want %q", input, got, expected)
	}
}

func TestRemoveZeroWidthChars(t *testing.T) {
	input := "hello\u200Bworld\uFEFF"
	expected := "helloworld"
	got := removeZeroWidthChars(input)
	if got != expected {
		t.Errorf("removeZeroWidthChars(%q) = %q, want %q", input, got, expected)
	}
}
