# Ghostty Theme TUI

A small Go TUI for browsing, previewing, and applying Ghostty themes.

## Usage

Run during development:

```bash
mise exec -- go run ./cmd/ghostty-theme-tui
```

Install locally from source:

```bash
mise exec -- make install GO="mise exec -- go"
```

By default this installs to:

```text
~/.local/bin/ghostty-theme-tui
```

Make sure `~/.local/bin` is in your `PATH`.

Install from a GitHub release without Go:

```bash
curl -fsSL https://raw.githubusercontent.com/Ryannnkl/ghostty-theme-tui/main/scripts/install.sh | sh
```

Optional color filter:

```bash
ghostty-theme-tui --color dark
ghostty-theme-tui --color light
ghostty-theme-tui --color all
```

## Controls

- Type to filter themes.
- Use `Up` and `Down` to move the selection.
- Press `Enter` to save the selected theme to the Ghostty config.
- Press `Esc` or `Ctrl+C` to exit without saving and restore the original preview.
- Press `/` to focus the search input.
- Press `Ctrl+R` to reload the theme list.

## Behavior

Themes are loaded with:

```bash
ghostty +list-themes --path
```

Preview is applied in the current terminal using OSC sequences for foreground,
background, cursor, and ANSI palette colors. The Ghostty config is only changed
after confirming with `Enter`.

The config path is:

```text
${XDG_CONFIG_HOME:-$HOME/.config}/ghostty/config
```

When saving, the tool writes `config.bak`, updates the first `theme = ...` line
or appends one if missing, then runs:

```bash
ghostty +validate-config
```

If validation fails, the previous config is restored.

## Releasing

Users do not need Go if you publish compiled release binaries.

Create and push a version tag:

```bash
git tag v0.1.0
git push origin v0.1.0
```

GitHub Actions will build release archives for:

- Linux amd64
- Linux arm64
- macOS amd64
- macOS arm64

After the release is published, users can install with:

```bash
curl -fsSL https://raw.githubusercontent.com/Ryannnkl/ghostty-theme-tui/main/scripts/install.sh | sh
```

To install a specific version:

```bash
curl -fsSL https://raw.githubusercontent.com/Ryannnkl/ghostty-theme-tui/main/scripts/install.sh | GHOSTTY_THEME_TUI_VERSION=v0.1.0 sh
```

To install somewhere else:

```bash
curl -fsSL https://raw.githubusercontent.com/Ryannnkl/ghostty-theme-tui/main/scripts/install.sh | GHOSTTY_THEME_TUI_BINDIR=/usr/local/bin sh
```
