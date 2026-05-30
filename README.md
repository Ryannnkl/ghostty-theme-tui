# Ghostty Theme TUI

A keyboard-driven terminal UI for browsing, previewing, and applying Ghostty
themes.

## Requirements

- Ghostty installed and available as `ghostty`
- Linux or macOS
- `curl` or `wget`

## Installation

Install the latest release:

```bash
curl -fsSL https://raw.githubusercontent.com/Ryannnkl/ghostty-theme-tui/main/scripts/install.sh | sh
```

This installs `ghostty-theme-tui` to:

```text
~/.local/bin/ghostty-theme-tui
```

If the command is not found after installing, add this to your shell config:

```bash
export PATH="$HOME/.local/bin:$PATH"
```

Then restart your shell.

## Usage

Open Ghostty and run:

```bash
ghostty-theme-tui
```

Filter themes by color:

```bash
ghostty-theme-tui --color dark
ghostty-theme-tui --color light
ghostty-theme-tui --color all
```

## Controls

- Type to search themes.
- Use `Up` and `Down` to move through the list.
- Press `Enter` to save the selected theme.
- Press `Esc` or `Ctrl+C` to exit without saving.
- Press `/` to focus search.
- Press `Ctrl+R` to reload themes.

## Updating

Run the installer again:

```bash
curl -fsSL https://raw.githubusercontent.com/Ryannnkl/ghostty-theme-tui/main/scripts/install.sh | sh
```

## Uninstalling

Remove the installed binary:

```bash
rm ~/.local/bin/ghostty-theme-tui
```

## What It Changes

The app previews themes in the current terminal while you browse. It only edits
your Ghostty config after you press `Enter`.

The config file is:

```text
${XDG_CONFIG_HOME:-$HOME/.config}/ghostty/config
```

Before saving, it creates:

```text
config.bak
```

If Ghostty rejects the updated config, the backup is restored automatically.
