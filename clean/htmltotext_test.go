package clean

import (
	"testing"
)

func TestHTMLToText_BasicList(t *testing.T) {
	html := `<ul><li>First item</li><li>Second item</li></ul>`
	expected := "- First item\n- Second item"
	got := HTMLToText(html)
	if got != expected {
		t.Errorf("HTMLToText basic list:\ngot:\n%s\n\nexpected:\n%s", got, expected)
	}
}

func TestHTMLToText_NestedList(t *testing.T) {
	html := `<ul>
		<li>Top item
			<ul>
				<li>Sub item 1</li>
				<li>Sub item 2</li>
			</ul>
		</li>
		<li>Another top</li>
	</ul>`
	got := HTMLToText(html)
	if !contains(got, "- Top item") {
		t.Errorf("Expected '- Top item' in:\n%s", got)
	}
	if !contains(got, "  - Sub item 1") {
		t.Errorf("Expected '  - Sub item 1' in:\n%s", got)
	}
	if !contains(got, "- Another top") {
		t.Errorf("Expected '- Another top' in:\n%s", got)
	}
}

func TestHTMLToText_OrderedList(t *testing.T) {
	html := `<ol><li>First</li><li>Second</li><li>Third</li></ol>`
	got := HTMLToText(html)
	if !contains(got, "1. First") {
		t.Errorf("Expected '1. First' in:\n%s", got)
	}
	if !contains(got, "2. Second") {
		t.Errorf("Expected '2. Second' in:\n%s", got)
	}
}

func TestHTMLToText_Links(t *testing.T) {
	html := `<p>Visit <a href="https://github.com/settings/billing">billing page</a> for details.</p>`
	got := HTMLToText(html)
	if !contains(got, "billing page (https://github.com/settings/billing)") {
		t.Errorf("Expected link with URL in:\n%s", got)
	}
}

func TestHTMLToText_LinkTextIsURL(t *testing.T) {
	html := `<a href="https://example.com">https://example.com</a>`
	got := HTMLToText(html)
	// Should NOT duplicate the URL
	if contains(got, "https://example.com (https://example.com)") {
		t.Errorf("URL should not be duplicated in:\n%s", got)
	}
	if !contains(got, "https://example.com") {
		t.Errorf("URL should be present in:\n%s", got)
	}
}

func TestHTMLToText_Headings(t *testing.T) {
	html := `<h1>Main Title</h1><p>Some text</p><h2>Subtitle</h2>`
	got := HTMLToText(html)
	if !contains(got, "Main Title") {
		t.Errorf("Expected 'Main Title' in:\n%s", got)
	}
	if !contains(got, "Some text") {
		t.Errorf("Expected 'Some text' in:\n%s", got)
	}
}

func TestHTMLToText_RemovesOpenInNewTab(t *testing.T) {
	html := `<a href="https://example.com">My Link (opens in a new tab)</a>`
	got := HTMLToText(html)
	if contains(got, "opens in a new tab") {
		t.Errorf("Should remove '(opens in a new tab)' from:\n%s", got)
	}
}

func TestHTMLToText_FullExample(t *testing.T) {
	html := `<div>
		<h2>1. How to visualize premium request consumption</h2>
		<h3>A. What individual users can see (self-service)</h3>
		<p>Each user can only see <strong>their own</strong> usage.</p>
		<p><strong>Where:</strong></p>
		<ul>
			<li>In the IDE (VS Code, JetBrains, Visual Studio, Xcode, etc.)
				<ul>
					<li>Copilot icon → quota / usage</li>
				</ul>
			</li>
			<li>On GitHub.com:
				<ul>
					<li><a href="https://github.com/settings/billing">https://github.com/settings/billing</a> (opens in a new tab)</li>
					<li>Scroll to <strong>Metered usage</strong></li>
					<li>Filter by <strong>Copilot</strong></li>
				</ul>
			</li>
		</ul>
		<p><strong>This shows:</strong></p>
		<ul>
			<li>How many premium requests they've used this month</li>
			<li>When the counter resets (always the 1st of the month, UTC)</li>
		</ul>
		<p>Source: <a href="https://docs.github.com/copilot">Monitoring your GitHub Copilot usage and entitlements (opens in a new tab)</a></p>
	</div>`

	got := HTMLToText(html)

	checks := []string{
		"1. How to visualize premium request consumption",
		"A. What individual users can see (self-service)",
		"their own",
		"- In the IDE",
		"  - Copilot icon",
		"- On GitHub.com:",
		"  - https://github.com/settings/billing",
		"  - Scroll to Metered usage",
		"  - Filter by Copilot",
		"- How many premium requests",
		"- When the counter resets",
		"Monitoring your GitHub Copilot usage and entitlements",
	}

	for _, check := range checks {
		if !contains(got, check) {
			t.Errorf("Expected to find %q in:\n%s", check, got)
		}
	}

	if contains(got, "opens in a new tab") {
		t.Errorf("Should not contain '(opens in a new tab)' in:\n%s", got)
	}
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && len(substr) > 0 && containsStr(s, substr)
}

func containsStr(s, sub string) bool {
	for i := 0; i <= len(s)-len(sub); i++ {
		if s[i:i+len(sub)] == sub {
			return true
		}
	}
	return false
}
