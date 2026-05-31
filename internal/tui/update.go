package tui

import (
	tea "github.com/charmbracelet/bubbletea"
)

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil

	case loadThemesMsg:
		m.loading = false
		m.themes = msg.themes
		m.currentTheme = msg.currentTheme
		m.originalTheme = msg.currentTheme
		m.applyFilter()
		m.selectCurrentTheme()
		if msg.err != nil {
			m.message = msg.err.Error()
			return m, nil
		}
		if theme, ok := m.selectedTheme(); ok {
			return m, applyPreview(m.previewer, theme)
		}
		return m, nil

	case saveThemeMsg:
		if msg.err != nil {
			m.message = msg.err.Error()
			return m, nil
		}
		if msg.reloadErr != nil {
			m.currentTheme = msg.name
			m.message = "Saved, but Ghostty reload failed. Press Ctrl+Shift+, to reload config."
			return m, nil
		}
		m.confirmed = true
		m.quitting = true
		m.message = "Saved " + msg.name
		return m, tea.Quit

	case previewMsg:
		if msg.err != nil {
			m.message = msg.err.Error()
		}
		return m, nil

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "esc":
			m.quitting = true
			if !m.confirmed {
				if theme, ok := m.originalThemeData(); ok {
					return m, tea.Batch(applyPreview(m.previewer, theme), tea.Quit)
				}
			}
			return m, tea.Quit

		case "enter":
			theme, ok := m.selectedTheme()
			if !ok {
				m.message = "No theme selected"
				return m, nil
			}
			m.message = "Saving " + theme.Name + "..."
			return m, saveTheme(theme.Name)

		case "up":
			if len(m.filtered) == 0 {
				return m, nil
			}
			if m.selected > 0 {
				m.selected--
			}
			theme, _ := m.selectedTheme()
			return m, applyPreview(m.previewer, theme)

		case "down":
			if len(m.filtered) == 0 {
				return m, nil
			}
			if m.selected < len(m.filtered)-1 {
				m.selected++
			}
			theme, _ := m.selectedTheme()
			return m, applyPreview(m.previewer, theme)

		case "/":
			m.input.Focus()
			return m, nil

		case "ctrl+r":
			m.loading = true
			m.message = "Reloading themes..."
			return m, loadThemes(m.color)
		}
	}

	before := m.input.Value()
	var cmd tea.Cmd
	m.input, cmd = m.input.Update(msg)
	if m.input.Value() != before {
		m.selected = 0
		m.applyFilter()
		if theme, ok := m.selectedTheme(); ok {
			return m, tea.Batch(cmd, applyPreview(m.previewer, theme))
		}
	}
	return m, cmd
}
