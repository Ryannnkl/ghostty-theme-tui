package preview

import (
	"bytes"
	"testing"

	"github.com/Ryannnkl/ghostty-theme-tui/internal/ghostty"
)

func TestApplyWritesOSCSequences(t *testing.T) {
	var out bytes.Buffer
	p := NewWithWriter(&out)

	err := p.Apply(ghostty.ThemeColors{
		Foreground: "#ffffff",
		Background: "#000000",
		Cursor:     "#ff00ff",
		Palette: map[int]string{
			1: "#111111",
			0: "#000000",
		},
	})
	if err != nil {
		t.Fatal(err)
	}

	want := "\x1b]10;#ffffff\x07\x1b]11;#000000\x07\x1b]12;#ff00ff\x07\x1b]4;0;#000000\x07\x1b]4;1;#111111\x07"
	if out.String() != want {
		t.Fatalf("unexpected output: %q", out.String())
	}
}
