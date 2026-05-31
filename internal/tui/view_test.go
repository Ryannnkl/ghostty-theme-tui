package tui

import (
	"fmt"
	"strings"
	"testing"

	"github.com/Ryannnkl/ghostty-theme-tui/internal/ghostty"
)

func TestRenderListKeepsSelectedThemeVisibleAtBottom(t *testing.T) {
	m := model{selected: 4}
	for i := 0; i < 10; i++ {
		m.filtered = append(m.filtered, ghostty.Theme{
			Name:   fmt.Sprintf("Theme %02d", i),
			Source: "resources",
		})
	}

	view := m.renderList(60, 5)

	if !strings.Contains(view, "Theme 04") {
		t.Fatalf("selected theme should be visible:\n%s", view)
	}
	if !strings.Contains(view, ">") {
		t.Fatalf("selected marker should be visible:\n%s", view)
	}
}
