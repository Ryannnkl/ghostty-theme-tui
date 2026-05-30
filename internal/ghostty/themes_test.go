package ghostty

import "testing"

func TestParseThemeList(t *testing.T) {
	input := []byte("Catppuccin Mocha (resources) /usr/share/ghostty/themes/Catppuccin Mocha\nMy Theme (user) /home/me/.config/ghostty/themes/My Theme\n")

	themes, err := ParseThemeList(input)
	if err != nil {
		t.Fatal(err)
	}

	if len(themes) != 2 {
		t.Fatalf("expected 2 themes, got %d", len(themes))
	}
	if themes[0].Name != "Catppuccin Mocha" {
		t.Fatalf("unexpected first name: %q", themes[0].Name)
	}
	if themes[0].Source != "resources" {
		t.Fatalf("unexpected first source: %q", themes[0].Source)
	}
	if themes[0].Path != "/usr/share/ghostty/themes/Catppuccin Mocha" {
		t.Fatalf("unexpected first path: %q", themes[0].Path)
	}
	if themes[1].Name != "My Theme" {
		t.Fatalf("unexpected second name: %q", themes[1].Name)
	}
}

func TestParseThemeColors(t *testing.T) {
	input := []byte(`
# comment
background = #1e1e2e
foreground = #cdd6f4
cursor-color = #f5e0dc
palette = 0=#45475a
palette = 1=#f38ba8
`)

	colors, err := ParseThemeColors(input)
	if err != nil {
		t.Fatal(err)
	}

	if colors.Background != "#1e1e2e" {
		t.Fatalf("unexpected background: %q", colors.Background)
	}
	if colors.Foreground != "#cdd6f4" {
		t.Fatalf("unexpected foreground: %q", colors.Foreground)
	}
	if colors.Cursor != "#f5e0dc" {
		t.Fatalf("unexpected cursor: %q", colors.Cursor)
	}
	if colors.Palette[0] != "#45475a" || colors.Palette[1] != "#f38ba8" {
		t.Fatalf("unexpected palette: %#v", colors.Palette)
	}
}
