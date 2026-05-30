package ghostty

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func ConfigPath() (string, error) {
	if xdg := os.Getenv("XDG_CONFIG_HOME"); xdg != "" {
		return filepath.Join(xdg, "ghostty", "config"), nil
	}
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, ".config", "ghostty", "config"), nil
}

func CurrentTheme() (string, error) {
	path, err := ConfigPath()
	if err != nil {
		return "", err
	}
	data, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		return "", nil
	}
	if err != nil {
		return "", err
	}
	return ParseCurrentTheme(data), nil
}

func ParseCurrentTheme(data []byte) string {
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
		if strings.TrimSpace(key) == "theme" {
			return strings.TrimSpace(value)
		}
	}
	return ""
}

func SaveTheme(name string) error {
	path, err := ConfigPath()
	if err != nil {
		return err
	}

	original, err := os.ReadFile(path)
	if err != nil && !os.IsNotExist(err) {
		return err
	}
	existed := err == nil

	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}

	backupPath := path + ".bak"
	if existed {
		if writeErr := os.WriteFile(backupPath, original, 0o644); writeErr != nil {
			return writeErr
		}
	}

	next := ReplaceTheme(original, name)
	if err := os.WriteFile(path, next, 0o644); err != nil {
		return err
	}

	if err := ValidateConfig(); err != nil {
		if restoreErr := restoreBackup(path, backupPath, original, existed); restoreErr != nil {
			return fmt.Errorf("%w; restore failed: %v", err, restoreErr)
		}
		return err
	}

	return nil
}

func ReplaceTheme(data []byte, name string) []byte {
	lines := strings.Split(string(data), "\n")
	replaced := false

	for i, line := range lines {
		trimmed := strings.TrimSpace(line)
		if trimmed == "" || strings.HasPrefix(trimmed, "#") {
			continue
		}
		key, _, ok := strings.Cut(line, "=")
		if ok && strings.TrimSpace(key) == "theme" {
			lines[i] = "theme = " + name
			replaced = true
			break
		}
	}

	if !replaced {
		if len(lines) == 1 && lines[0] == "" {
			lines[0] = "theme = " + name
		} else {
			if lines[len(lines)-1] != "" {
				lines = append(lines, "")
			}
			lines = append(lines, "theme = "+name)
		}
	}

	return []byte(strings.Join(lines, "\n"))
}

func ValidateConfig() error {
	cmd := exec.Command("ghostty", "+validate-config")
	if output, err := cmd.CombinedOutput(); err != nil {
		message := strings.TrimSpace(string(output))
		if message == "" {
			return fmt.Errorf("validate config: %w", err)
		}
		return fmt.Errorf("validate config: %s", message)
	}
	return nil
}

func restoreBackup(path, backupPath string, original []byte, existed bool) error {
	if !existed {
		return os.Remove(path)
	}
	if _, err := os.Stat(backupPath); err == nil {
		data, readErr := os.ReadFile(backupPath)
		if readErr != nil {
			return readErr
		}
		return os.WriteFile(path, data, 0o644)
	}
	return os.WriteFile(path, original, 0o644)
}
