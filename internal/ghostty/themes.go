package ghostty

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
)

var themeLinePattern = regexp.MustCompile(`^(.+) \((resources|user)\) (.+)$`)

type Theme struct {
	Name   string
	Source string
	Path   string
	Colors ThemeColors
}

type ThemeColors struct {
	Foreground string
	Background string
	Cursor     string
	Palette    map[int]string
}

func ListThemes(color string) ([]Theme, error) {
	args := []string{"+list-themes", "--path"}
	if color == "dark" || color == "light" {
		args = []string{"+list-themes", "--color=" + color, "--path"}
	}

	cmd := exec.Command("ghostty", args...)
	output, err := cmd.Output()
	if err != nil {
		return nil, commandError("list themes", err)
	}

	themes, err := ParseThemeList(output)
	if err != nil {
		return nil, err
	}

	for i := range themes {
		colors, err := ParseThemeFile(themes[i].Path)
		if err == nil {
			themes[i].Colors = colors
		}
	}

	return themes, nil
}

func ParseThemeList(data []byte) ([]Theme, error) {
	scanner := bufio.NewScanner(bytes.NewReader(data))
	var themes []Theme
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}

		matches := themeLinePattern.FindStringSubmatch(line)
		if matches == nil {
			return nil, fmt.Errorf("invalid theme list line: %q", line)
		}

		themes = append(themes, Theme{
			Name:   strings.TrimSpace(matches[1]),
			Source: matches[2],
			Path:   strings.TrimSpace(matches[3]),
		})
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return themes, nil
}

func ParseThemeFile(path string) (ThemeColors, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return ThemeColors{}, err
	}
	return ParseThemeColors(data)
}

func ParseThemeColors(data []byte) (ThemeColors, error) {
	colors := ThemeColors{Palette: make(map[int]string)}
	scanner := bufio.NewScanner(bytes.NewReader(data))

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		key, value, ok := strings.Cut(line, "=")
		if !ok {
			continue
		}

		key = strings.TrimSpace(key)
		value = strings.TrimSpace(value)

		switch key {
		case "foreground":
			colors.Foreground = value
		case "background":
			colors.Background = value
		case "cursor-color":
			colors.Cursor = value
		case "palette":
			indexText, color, ok := strings.Cut(value, "=")
			if !ok {
				continue
			}
			index, err := strconv.Atoi(strings.TrimSpace(indexText))
			if err != nil {
				continue
			}
			if index < 0 || index > 255 {
				continue
			}
			colors.Palette[index] = strings.TrimSpace(color)
		}
	}

	if err := scanner.Err(); err != nil {
		return ThemeColors{}, err
	}
	return colors, nil
}

func commandError(action string, err error) error {
	var exitErr *exec.ExitError
	if errors.As(err, &exitErr) {
		stderr := strings.TrimSpace(string(exitErr.Stderr))
		if stderr != "" {
			return fmt.Errorf("%s: %s", action, stderr)
		}
	}
	return fmt.Errorf("%s: %w", action, err)
}
