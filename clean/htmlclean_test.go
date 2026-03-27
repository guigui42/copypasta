package clean

import (
	"strings"
	"testing"
)

func TestCleanHTML_StripsDarkBackground(t *testing.T) {
	input := `<div style="background-color: #1a1a2e; color: #e0e0e0; font-family: monospace;">
		<h2 style="color: #7ecbf5; background: #2d2d44;">Title</h2>
		<p style="background-color: rgb(30,30,50); color: white;">Some text</p>
	</div>`

	got := CleanHTML(input)

	if strings.Contains(got, "background-color") {
		t.Errorf("Should strip background-color, got:\n%s", got)
	}
	if strings.Contains(got, "#1a1a2e") {
		t.Errorf("Should strip dark color values, got:\n%s", got)
	}
	if strings.Contains(got, "font-family") {
		t.Errorf("Should strip font-family, got:\n%s", got)
	}
	if !strings.Contains(got, "Title") {
		t.Errorf("Should preserve text content, got:\n%s", got)
	}
	if !strings.Contains(got, "Some text") {
		t.Errorf("Should preserve text content, got:\n%s", got)
	}
}

func TestCleanHTML_PreservesListStructure(t *testing.T) {
	input := `<ul style="color: #ccc; background: #222;">
		<li style="color: white;">First item</li>
		<li>Second item
			<ul><li>Nested</li></ul>
		</li>
	</ul>`

	got := CleanHTML(input)

	if !strings.Contains(got, "<ul") {
		t.Errorf("Should preserve <ul> tags, got:\n%s", got)
	}
	if !strings.Contains(got, "<li>") || !strings.Contains(got, "<li") {
		t.Errorf("Should preserve <li> tags, got:\n%s", got)
	}
	if !strings.Contains(got, "First item") {
		t.Errorf("Should preserve text, got:\n%s", got)
	}
	if !strings.Contains(got, "Nested") {
		t.Errorf("Should preserve nested text, got:\n%s", got)
	}
}

func TestCleanHTML_PreservesLinks(t *testing.T) {
	input := `<a href="https://github.com/settings/billing" style="color: #58a6ff;">billing page</a>`

	got := CleanHTML(input)

	if !strings.Contains(got, `href="https://github.com/settings/billing"`) {
		t.Errorf("Should preserve href, got:\n%s", got)
	}
	if !strings.Contains(got, "billing page") {
		t.Errorf("Should preserve link text, got:\n%s", got)
	}
	if strings.Contains(got, "#58a6ff") {
		t.Errorf("Should strip color from link style, got:\n%s", got)
	}
}

func TestCleanHTML_RemovesOpenInNewTab(t *testing.T) {
	input := `<a href="https://example.com">My Link (opens in a new tab)</a>`
	got := CleanHTML(input)
	if strings.Contains(got, "opens in a new tab") {
		t.Errorf("Should remove '(opens in a new tab)', got:\n%s", got)
	}
}

func TestCleanHTML_StripsStyleScriptTags(t *testing.T) {
	input := `<div><style>.dark { color: white; }</style><p>Keep this</p><script>alert('x')</script></div>`
	got := CleanHTML(input)
	if strings.Contains(got, "<style") {
		t.Errorf("Should remove <style>, got:\n%s", got)
	}
	if strings.Contains(got, "<script") {
		t.Errorf("Should remove <script>, got:\n%s", got)
	}
	if !strings.Contains(got, "Keep this") {
		t.Errorf("Should keep text content, got:\n%s", got)
	}
}

func TestCleanHTML_StripsClassAndId(t *testing.T) {
	input := `<div class="dark-theme" id="main"><p class="text-white">Hello</p></div>`
	got := CleanHTML(input)
	if strings.Contains(got, `class=`) {
		t.Errorf("Should strip class attributes, got:\n%s", got)
	}
	if strings.Contains(got, `id=`) {
		t.Errorf("Should strip id attributes, got:\n%s", got)
	}
}
