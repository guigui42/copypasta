package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/guigui42/copypasta/clean"
	"github.com/guigui42/copypasta/clipboard"
)

func main() {
	plainText, _ := clipboard.Read()
	htmlText, _ := clipboard.ReadHTML()

	if strings.TrimSpace(plainText) == "" && strings.TrimSpace(htmlText) == "" {
		fmt.Fprintln(os.Stderr, "Clipboard is empty — copy some text first!")
		os.Exit(1)
	}

	if strings.TrimSpace(htmlText) != "" {
		cleanedHTML := clean.CleanHTML(htmlText)
		plainFallback := clean.HTMLToText(htmlText)
		if err := clipboard.WriteHTML(cleanedHTML, plainFallback); err != nil {
			fmt.Fprintf(os.Stderr, "Error writing clipboard: %v\n", err)
			os.Exit(1)
		}
	} else {
		cleaned := clean.Clean(plainText)
		if err := clipboard.Write(cleaned); err != nil {
			fmt.Fprintf(os.Stderr, "Error writing clipboard: %v\n", err)
			os.Exit(1)
		}
	}

	fmt.Println("✓ Cleaned text copied to clipboard!")
}
