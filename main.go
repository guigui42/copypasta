package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/guigui42/copypasta/clean"
	"github.com/guigui42/copypasta/clipboard"
)

func main() {
	if len(os.Args) > 1 && os.Args[1] == "--install" {
		if err := installQuickAction(); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		return
	}

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

func installQuickAction() error {
	// Resolve the absolute path of the running binary
	exe, err := os.Executable()
	if err != nil {
		return fmt.Errorf("could not determine executable path: %w", err)
	}
	exe, err = filepath.EvalSymlinks(exe)
	if err != nil {
		return fmt.Errorf("could not resolve symlinks: %w", err)
	}

	home, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("could not determine home directory: %w", err)
	}

	workflowDir := filepath.Join(home, "Library", "Services", "Clean Clipboard.workflow", "Contents")
	if err := os.MkdirAll(workflowDir, 0755); err != nil {
		return fmt.Errorf("could not create workflow directory: %w", err)
	}

	infoPlist := `<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
<dict>
	<key>NSServices</key>
	<array>
		<dict>
			<key>NSMenuItem</key>
			<dict>
				<key>default</key>
				<string>Clean Clipboard</string>
			</dict>
			<key>NSMessage</key>
			<string>runWorkflowAsService</string>
		</dict>
	</array>
</dict>
</plist>`

	command := exe + " 2>&1 || true"
	documentWflow := `<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
<dict>
	<key>AMApplicationBuild</key>
	<string>523</string>
	<key>AMApplicationVersion</key>
	<string>2.10</string>
	<key>AMDocumentVersion</key>
	<string>2</string>
	<key>actions</key>
	<array>
		<dict>
			<key>action</key>
			<dict>
				<key>AMAccepts</key>
				<dict>
					<key>Container</key>
					<string>List</string>
					<key>Optional</key>
					<true/>
					<key>Types</key>
					<array>
						<string>com.apple.cocoa.string</string>
					</array>
				</dict>
				<key>AMActionVersion</key>
				<string>2.0.3</string>
				<key>AMApplication</key>
				<array>
					<string>Automator</string>
				</array>
				<key>AMCategory</key>
				<string>AMCategoryUtilities</string>
				<key>AMIconName</key>
				<string>RunShellScript</string>
				<key>AMParameterProperties</key>
				<dict>
					<key>COMMAND_STRING</key>
					<dict/>
					<key>inputMethod</key>
					<dict/>
					<key>shell</key>
					<dict/>
					<key>source</key>
					<dict/>
				</dict>
				<key>AMProvides</key>
				<dict>
					<key>Container</key>
					<string>List</string>
					<key>Types</key>
					<array>
						<string>com.apple.cocoa.string</string>
					</array>
				</dict>
				<key>AMRequiredResources</key>
				<array/>
				<key>ActionBundlePath</key>
				<string>/System/Library/Automator/Run Shell Script.action</string>
				<key>ActionName</key>
				<string>Run Shell Script</string>
				<key>ActionParameters</key>
				<dict>
					<key>COMMAND_STRING</key>
					<string>` + command + `</string>
					<key>CheckedForUserDefaultShell</key>
					<true/>
					<key>inputMethod</key>
					<integer>1</integer>
					<key>shell</key>
					<string>/bin/zsh</string>
					<key>source</key>
					<string></string>
				</dict>
				<key>BundleIdentifier</key>
				<string>com.apple.RunShellScript</string>
				<key>CFBundleVersion</key>
				<string>2.0.3</string>
				<key>CanShowSelectedItemsWhenRun</key>
				<false/>
				<key>CanShowWhenRun</key>
				<true/>
				<key>GroupedWorkflow</key>
				<true/>
			</dict>
			<key>isViewVisible</key>
			<true/>
		</dict>
	</array>
	<key>connectors</key>
	<dict/>
	<key>workflowMetaData</key>
	<dict>
		<key>workflowTypeIdentifier</key>
		<string>com.apple.Automator.servicesMenu</string>
	</dict>
</dict>
</plist>`

	if err := os.WriteFile(filepath.Join(workflowDir, "Info.plist"), []byte(infoPlist), 0644); err != nil {
		return fmt.Errorf("could not write Info.plist: %w", err)
	}
	if err := os.WriteFile(filepath.Join(workflowDir, "document.wflow"), []byte(documentWflow), 0644); err != nil {
		return fmt.Errorf("could not write document.wflow: %w", err)
	}

	fmt.Println("✓ Quick Action installed: ~/Library/Services/Clean Clipboard.workflow")
	fmt.Println()
	fmt.Println("Next step — assign a keyboard shortcut:")
	fmt.Println("  System Settings → Keyboard → Keyboard Shortcuts → Services → General")
	fmt.Println("  Find \"Clean Clipboard\" and set your shortcut (e.g. ⌃⌥⌘V)")
	return nil
}
