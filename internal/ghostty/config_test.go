package ghostty

import "testing"

func TestParseCurrentTheme(t *testing.T) {
	input := []byte(`
# theme = Ignored
font-family = JetBrains Mono
theme = Catppuccin Mocha
`)

	got := ParseCurrentTheme(input)
	if got != "Catppuccin Mocha" {
		t.Fatalf("expected current theme, got %q", got)
	}
}

func TestReplaceThemeReplacesFirstThemeLine(t *testing.T) {
	input := []byte("font-size = 14\ntheme = Old Theme\nwindow-padding-x = 8\n")

	got := string(ReplaceTheme(input, "New Theme"))
	want := "font-size = 14\ntheme = New Theme\nwindow-padding-x = 8\n"
	if got != want {
		t.Fatalf("unexpected config:\n%s", got)
	}
}

func TestReplaceThemeAppendsWhenMissing(t *testing.T) {
	input := []byte("font-size = 14\n")

	got := string(ReplaceTheme(input, "New Theme"))
	want := "font-size = 14\n\ntheme = New Theme"
	if got != want {
		t.Fatalf("unexpected config:\n%s", got)
	}
}
