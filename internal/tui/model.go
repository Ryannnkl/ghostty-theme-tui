package tui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"

	"github.com/Ryannnkl/ghostty-theme-tui/internal/ghostty"
	"github.com/Ryannnkl/ghostty-theme-tui/internal/preview"
)

type model struct {
	themes        []ghostty.Theme
	filtered      []ghostty.Theme
	selected      int
	input         textinput.Model
	currentTheme  string
	originalTheme string
	message       string
	width         int
	height        int
	color         string
	loading       bool
	quitting      bool
	confirmed     bool
	previewer     preview.Previewer
}

type loadThemesMsg struct {
	themes       []ghostty.Theme
	currentTheme string
	err          error
}

type saveThemeMsg struct {
	name      string
	err       error
	reloadErr error
}

type previewMsg struct {
	err error
}

func New(color string) tea.Model {
	input := textinput.New()
	input.Placeholder = "Search themes"
	input.Prompt = "/ "
	input.Focus()
	input.CharLimit = 120

	return model{
		input:     input,
		color:     color,
		loading:   true,
		previewer: preview.New(),
	}
}

func (m model) Init() tea.Cmd {
	return loadThemes(m.color)
}

func loadThemes(color string) tea.Cmd {
	return func() tea.Msg {
		currentTheme, currentErr := ghostty.CurrentTheme()
		themes, themeErr := ghostty.ListThemes(color)
		if themeErr != nil {
			return loadThemesMsg{currentTheme: currentTheme, err: themeErr}
		}
		if currentErr != nil {
			return loadThemesMsg{themes: themes, currentTheme: currentTheme, err: currentErr}
		}
		return loadThemesMsg{themes: themes, currentTheme: currentTheme}
	}
}

func saveTheme(name string) tea.Cmd {
	return func() tea.Msg {
		if err := ghostty.SaveTheme(name); err != nil {
			return saveThemeMsg{name: name, err: err}
		}
		_, reloadErr := ghostty.ReloadConfig()
		return saveThemeMsg{name: name, reloadErr: reloadErr}
	}
}

func applyPreview(p preview.Previewer, theme ghostty.Theme) tea.Cmd {
	return func() tea.Msg {
		if err := p.Apply(theme.Colors); err != nil {
			return previewMsg{err: err}
		}
		return previewMsg{}
	}
}

func (m model) selectedTheme() (ghostty.Theme, bool) {
	if m.selected < 0 || m.selected >= len(m.filtered) {
		return ghostty.Theme{}, false
	}
	return m.filtered[m.selected], true
}

func (m model) originalThemeData() (ghostty.Theme, bool) {
	if m.originalTheme == "" {
		return ghostty.Theme{}, false
	}
	for _, theme := range m.themes {
		if strings.EqualFold(theme.Name, m.originalTheme) {
			return theme, true
		}
	}
	return ghostty.Theme{}, false
}

func (m *model) applyFilter() {
	query := strings.ToLower(strings.TrimSpace(m.input.Value()))
	m.filtered = m.filtered[:0]
	for _, theme := range m.themes {
		if query == "" || strings.Contains(strings.ToLower(theme.Name), query) {
			m.filtered = append(m.filtered, theme)
		}
	}

	if len(m.filtered) == 0 {
		m.selected = 0
		return
	}
	if m.selected >= len(m.filtered) {
		m.selected = len(m.filtered) - 1
	}
	if m.selected < 0 {
		m.selected = 0
	}
}

func (m *model) selectCurrentTheme() {
	if m.currentTheme == "" {
		m.selected = 0
		return
	}
	for i, theme := range m.filtered {
		if strings.EqualFold(theme.Name, m.currentTheme) {
			m.selected = i
			return
		}
	}
	m.selected = 0
}

func (m model) footerMessage() string {
	if m.message != "" {
		return m.message
	}
	if m.loading {
		return "Loading themes..."
	}
	if len(m.filtered) == 0 {
		return "No themes found"
	}
	return fmt.Sprintf("%d/%d themes", len(m.filtered), len(m.themes))
}
