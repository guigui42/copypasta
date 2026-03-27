package clean

import (
	"fmt"
	"net/url"
	"strings"

	"golang.org/x/net/html"
)

// HTMLToText converts HTML to well-structured plain text, preserving
// list hierarchy, indentation, links, and headings.
func HTMLToText(htmlStr string) string {
	doc, err := html.Parse(strings.NewReader(htmlStr))
	if err != nil {
		return htmlStr
	}

	ctx := &htmlCtx{}
	ctx.walk(doc)

	result := ctx.buf.String()

	// Post-process: clean up excessive blank lines and trailing whitespace
	result = collapseBlankLines(result)
	result = trimTrailingWhitespace(result)
	result = removeOpenInNewTab(result)
	result = removeZeroWidthChars(result)
	result = strings.TrimSpace(result)

	return result
}

type listInfo struct {
	ordered bool
	index   int
}

type htmlCtx struct {
	buf        strings.Builder
	listStack  []listInfo
	inPre      bool
	lastWasNL  bool
	inAnchor   bool
	anchorHref string
	anchorText strings.Builder
}

func (c *htmlCtx) write(s string) {
	if s == "" {
		return
	}
	c.buf.WriteString(s)
	c.lastWasNL = strings.HasSuffix(s, "\n")
}

func (c *htmlCtx) ensureNewline() {
	if !c.lastWasNL {
		c.buf.WriteString("\n")
		c.lastWasNL = true
	}
}

func (c *htmlCtx) ensureBlankLine() {
	c.ensureNewline()
	s := c.buf.String()
	if !strings.HasSuffix(s, "\n\n") {
		c.buf.WriteString("\n")
	}
}

func (c *htmlCtx) listDepth() int {
	return len(c.listStack)
}

func (c *htmlCtx) indent() string {
	depth := c.listDepth()
	if depth <= 1 {
		return ""
	}
	return strings.Repeat("  ", depth-1)
}

func (c *htmlCtx) walk(n *html.Node) {
	switch n.Type {
	case html.TextNode:
		text := n.Data
		if !c.inPre {
			// Collapse whitespace in non-pre blocks
			text = collapseWhitespace(text)
		}
		if c.inAnchor {
			c.anchorText.WriteString(text)
		} else {
			c.write(text)
		}
		return

	case html.ElementNode:
		c.handleElement(n)
		return
	}

	for child := n.FirstChild; child != nil; child = child.NextSibling {
		c.walk(child)
	}
}

func (c *htmlCtx) handleElement(n *html.Node) {
	tag := strings.ToLower(n.Data)

	switch tag {
	case "h1", "h2", "h3", "h4", "h5", "h6":
		c.ensureBlankLine()
		c.walkChildren(n)
		c.ensureNewline()
		c.write("\n")

	case "p":
		c.ensureBlankLine()
		c.walkChildren(n)
		c.ensureNewline()

	case "br":
		c.ensureNewline()

	case "ul":
		if c.listDepth() == 0 {
			c.ensureBlankLine()
		}
		c.listStack = append(c.listStack, listInfo{ordered: false})
		c.walkChildren(n)
		c.listStack = c.listStack[:len(c.listStack)-1]
		if c.listDepth() == 0 {
			c.ensureNewline()
		}

	case "ol":
		if c.listDepth() == 0 {
			c.ensureBlankLine()
		}
		c.listStack = append(c.listStack, listInfo{ordered: true, index: 0})
		c.walkChildren(n)
		c.listStack = c.listStack[:len(c.listStack)-1]
		if c.listDepth() == 0 {
			c.ensureNewline()
		}

	case "li":
		c.ensureNewline()
		indent := c.indent()
		if len(c.listStack) > 0 {
			li := &c.listStack[len(c.listStack)-1]
			if li.ordered {
				li.index++
				c.write(fmt.Sprintf("%s%d. ", indent, li.index))
			} else {
				c.write(indent + "- ")
			}
		}
		c.walkChildren(n)

	case "a":
		href := getAttr(n, "href")
		c.inAnchor = true
		c.anchorHref = href
		c.anchorText.Reset()
		c.walkChildren(n)
		linkText := strings.TrimSpace(c.anchorText.String())
		c.inAnchor = false

		// If the visible text IS the URL, just write it once
		if linkText == href || linkText == "" {
			c.write(href)
		} else if isValidURL(href) && href != "" {
			c.write(fmt.Sprintf("%s (%s)", linkText, href))
		} else {
			c.write(linkText)
		}

	case "pre", "code":
		if tag == "pre" {
			c.inPre = true
		}
		c.walkChildren(n)
		if tag == "pre" {
			c.inPre = false
		}

	case "div":
		c.ensureNewline()
		c.walkChildren(n)
		c.ensureNewline()

	case "hr":
		c.ensureBlankLine()
		c.write("---\n")

	case "style", "script", "noscript":
		// Skip these entirely

	case "strong", "b", "em", "i", "u", "mark", "span", "sup", "sub":
		c.walkChildren(n)

	default:
		c.walkChildren(n)
	}
}

func (c *htmlCtx) walkChildren(n *html.Node) {
	for child := n.FirstChild; child != nil; child = child.NextSibling {
		c.walk(child)
	}
}

func getAttr(n *html.Node, key string) string {
	for _, a := range n.Attr {
		if a.Key == key {
			return a.Val
		}
	}
	return ""
}

func isValidURL(s string) bool {
	u, err := url.Parse(s)
	return err == nil && (u.Scheme == "http" || u.Scheme == "https")
}

func collapseWhitespace(s string) string {
	var buf strings.Builder
	prevSpace := false
	for _, r := range s {
		if r == ' ' || r == '\t' || r == '\n' || r == '\r' {
			if !prevSpace {
				buf.WriteRune(' ')
			}
			prevSpace = true
		} else {
			buf.WriteRune(r)
			prevSpace = false
		}
	}
	return buf.String()
}
