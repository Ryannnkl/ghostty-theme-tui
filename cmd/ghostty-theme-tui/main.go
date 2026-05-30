package main

import (
	"flag"
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/Ryannnkl/ghostty-theme-tui/internal/tui"
)

func main() {
	color := flag.String("color", "all", "theme color filter: dark, light, or all")
	flag.Parse()

	if *color != "all" && *color != "dark" && *color != "light" {
		fmt.Fprintln(os.Stderr, "--color must be one of: all, dark, light")
		os.Exit(2)
	}

	model := tui.New(*color)
	program := tea.NewProgram(model, tea.WithAltScreen())
	if _, err := program.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "ghostty-theme-tui: %v\n", err)
		os.Exit(1)
	}
}
