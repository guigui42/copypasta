package clean

import (
	"regexp"
	"strings"

	"golang.org/x/net/html"
)

// CleanHTML takes raw HTML from the clipboard and returns sanitized HTML
// that preserves structure (lists, links, headings, bold, italic) but
// strips dark backgrounds, colors, custom fonts, and emoji noise.
func CleanHTML(htmlStr string) string {
	doc, err := html.Parse(strings.NewReader(htmlStr))
	if err != nil {
		return htmlStr
	}

	sanitizeNode(doc)

	var buf strings.Builder
	html.Render(&buf, doc)

	result := buf.String()

	// Remove "(opens in a new tab)" text
	result = openInNewTabRe.ReplaceAllString(result, "")

	return result
}

// sanitizeNode recursively cleans the HTML tree.
func sanitizeNode(n *html.Node) {
	// Remove nodes that should be stripped entirely
	if n.Type == html.ElementNode {
		tag := strings.ToLower(n.Data)
		if tag == "style" || tag == "script" || tag == "noscript" || tag == "svg" || tag == "hr" {
			// Mark for removal
			n.Data = "removed"
			n.FirstChild = nil
			n.LastChild = nil
			return
		}
	}

	// Strip text nodes that are just emoji
	if n.Type == html.TextNode {
		n.Data = stripStandaloneEmoji(n.Data)
	}

	// Clean style attributes on elements
	if n.Type == html.ElementNode {
		cleanAttributes(n)
	}

	// Recurse, collecting children to remove
	var toRemove []*html.Node
	for child := n.FirstChild; child != nil; child = child.NextSibling {
		sanitizeNode(child)
		if child.Type == html.ElementNode && child.Data == "removed" {
			toRemove = append(toRemove, child)
		}
	}
	for _, child := range toRemove {
		n.RemoveChild(child)
	}
}

// cleanAttributes strips style properties related to colors/backgrounds
// and removes class attributes, keeping only structural attributes.
func cleanAttributes(n *html.Node) {
	var cleaned []html.Attribute
	for _, attr := range n.Attr {
		switch attr.Key {
		case "style":
			// Strip color/background properties but keep layout ones
			sanitized := sanitizeStyle(attr.Val)
			if sanitized != "" {
				cleaned = append(cleaned, html.Attribute{Key: "style", Val: sanitized})
			}
		case "href":
			// Keep links
			cleaned = append(cleaned, attr)
		case "src", "alt":
			// Keep image basics
			cleaned = append(cleaned, attr)
		case "colspan", "rowspan":
			// Keep table structure
			cleaned = append(cleaned, attr)
			// Drop: class, id, data-*, bgcolor, color, width, height on non-img, etc.
		}
	}
	n.Attr = cleaned
}

// Properties to strip from inline styles (colors, backgrounds, fonts, etc.)
var stripStyleProps = regexp.MustCompile(
	`(?i)(background[\w-]*|color|font-family|font-size|text-shadow|box-shadow|border[\w-]*|outline[\w-]*)\s*:[^;]*;?`,
)

func sanitizeStyle(style string) string {
	cleaned := stripStyleProps.ReplaceAllString(style, "")
	cleaned = strings.TrimSpace(cleaned)
	// Remove trailing semicolons and whitespace
	cleaned = strings.TrimRight(cleaned, "; ")
	if cleaned == "" {
		return ""
	}
	return cleaned
}

// stripStandaloneEmoji removes emoji characters that appear alone (not within text).
var emojiRe = regexp.MustCompile(`^[\s]*[\x{1F300}-\x{1FAF8}\x{2600}-\x{27BF}\x{FE00}-\x{FE0F}\x{200D}\x{20E3}\x{E0020}-\x{E007F}✅❌⚠📷🔗💡🎯🚀✓✔☑☐📌📝🔒🔑📘]+[\s]*$`)

func stripStandaloneEmoji(s string) string {
	if emojiRe.MatchString(s) {
		return ""
	}
	return s
}
