package clipboard

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"
)

// Read returns the current clipboard contents as plain text.
func Read() (string, error) {
	cmd := exec.Command("pbpaste")
	var out bytes.Buffer
	cmd.Stdout = &out
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("failed to read clipboard: %w", err)
	}
	return out.String(), nil
}

// ReadHTML returns the HTML content from the clipboard if available.
// When copying from a browser, the clipboard contains both plain text and HTML.
// The HTML preserves structure (lists, links, headings) that plain text loses.
func ReadHTML() (string, error) {
	script := `ObjC.import("AppKit");
var pb = $.NSPasteboard.generalPasteboard;
var html = pb.stringForType("public.html");
html ? html.js : "";`
	cmd := exec.Command("osascript", "-l", "JavaScript", "-e", script)
	var out bytes.Buffer
	cmd.Stdout = &out
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("failed to read HTML clipboard: %w", err)
	}
	return strings.TrimSpace(out.String()), nil
}

// Write copies the given text to the clipboard as plain text.
func Write(text string) error {
	cmd := exec.Command("pbcopy")
	cmd.Stdin = strings.NewReader(text)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to write clipboard: %w", err)
	}
	return nil
}

// WriteHTML copies cleaned HTML to the clipboard so that rich-text apps
// (Outlook, Word, etc.) receive it with formatting intact.
// It sets both HTML and plain text on the clipboard.
func WriteHTML(htmlContent string, plainFallback string) error {
	// Use osascript to set HTML and plain text on the clipboard.
	// We set NSPasteboardTypeHTML (which also registers "Apple HTML pasteboard type")
	// because Outlook/Word look for that type specifically.
	script := fmt.Sprintf(`ObjC.import("AppKit");
var pb = $.NSPasteboard.generalPasteboard;
pb.clearContents;
pb.setStringForType($("%s"), $.NSPasteboardTypeHTML);
pb.setStringForType($("%s"), "public.utf8-plain-text");
"ok";`,
		jsEscape(htmlContent),
		jsEscape(plainFallback),
	)
	cmd := exec.Command("osascript", "-l", "JavaScript", "-e", script)
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to write HTML clipboard: %w (%s)", err, stderr.String())
	}
	return nil
}

// jsEscape escapes a string for embedding in a JavaScript string literal.
func jsEscape(s string) string {
	r := strings.NewReplacer(
		`\`, `\\`,
		`"`, `\"`,
		"\n", `\n`,
		"\r", `\r`,
		"\t", `\t`,
	)
	return r.Replace(s)
}
