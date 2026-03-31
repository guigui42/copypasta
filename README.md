# 🧹 copypasta

A tiny Go tool that cleans up text copied from dark-themed websites so it pastes cleanly into Outlook, Word, and other rich-text editors.

## The Problem

When you copy text from sites with dark backgrounds (like GitHub Copilot Support), pasting into Outlook/Word preserves the dark background, colors, and formatting noise — making it unreadable.

## How It Works

1. Copy text from any website
2. Run `copypasta`
3. Paste into Outlook, Word, etc. — clean and readable!

The tool reads your clipboard (including the hidden HTML), strips dark backgrounds, colors, borders, and noise while **keeping** lists, links, headings, and bold/italic formatting. It writes both clean HTML and a plain-text fallback back to your clipboard.

## What It Cleans

| Stripped | Kept |
|----------|------|
| Background colors | Bullet & numbered lists |
| Text colors & fonts | Clickable links with URLs |
| Borders & shadows | Headings & bold/italic |
| `(opens in a new tab)` annotations | Paragraph structure |
| Horizontal rules | Tables |
| Emoji prefixes | Indentation |
| Zero-width characters | |
| `<style>` / `<script>` tags | |

## Installation

### Homebrew (recommended)

```bash
brew install guigui42/tap/copypasta
```

### From source

```bash
# Clone and build
git clone https://github.com/guigui42/copypasta.git
cd copypasta
go build -o copypasta .

# Or install directly
go install github.com/guigui42/copypasta@latest
```

## Usage

```bash
copypasta
# ✓ Cleaned text copied to clipboard!
```

That's it. One command, no flags, no UI.

## macOS Quick Action (no third-party app)

Create a system-wide keyboard shortcut using macOS Automator — no Raycast or other apps needed:

1. Install the Quick Action:

```bash
copypasta --install
```

2. Assign a keyboard shortcut:
   - **System Settings → Keyboard → Keyboard Shortcuts → Services → General**
   - Find **Clean Clipboard**, click "Add Shortcut", press your combo (e.g. `⌃⌥⌘V`)

**Workflow:** Copy → press hotkey → Paste clean.

> **Note:** The first time you use the shortcut, macOS may ask for Accessibility permissions — grant them.

## Raycast Integration

For a seamless one-keystroke workflow, add a Raycast script command:

1. Create `~/Documents/raycast-scripts/copypasta.sh`:

```bash
#!/bin/bash

# @raycast.schemaVersion 1
# @raycast.title Clean Clipboard
# @raycast.mode silent
# @raycast.icon 🧹
# @raycast.packageName Copypasta

/path/to/copypasta
```

2. `chmod +x ~/Documents/raycast-scripts/copypasta.sh`
3. Raycast → Settings → Extensions → Script Commands → Add Script Directory → select `~/Documents/raycast-scripts`
4. Assign a hotkey (e.g. `⌃⌥⌘V`)

**Workflow:** Copy → press hotkey → Paste clean.

## License

MIT
