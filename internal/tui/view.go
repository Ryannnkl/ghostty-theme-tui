package tui

import (
	"fmt"
	"sort"
	"strings"

	"github.com/charmbracelet/lipgloss"

	"github.com/Ryannnkl/ghostty-theme-tui/internal/ghostty"
)

const (
	minWidth       = 40
	previewAt      = 92
	previewWidth   = 34
	pageHPadding   = 2
	listMinHeight  = 7
	sourceColWidth = 9
	stateColWidth  = 9
)

var (
	appStyle = lipgloss.NewStyle().Padding(1, pageHPadding)

	titleStyle = lipgloss.NewStyle().Bold(true)

	headerMetaStyle = lipgloss.NewStyle().Faint(true)

	searchStyle = lipgloss.NewStyle().
			Border(lipgloss.NormalBorder(), false, false, true, false).
			BorderForeground(lipgloss.Color("8")).
			PaddingBottom(1)

	selectedStyle = lipgloss.NewStyle().
			Bold(true).
			Reverse(true)

	currentStyle = lipgloss.NewStyle().Bold(true)

	mutedStyle = lipgloss.NewStyle().Faint(true)

	errorStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("9")).
			Bold(true)

	previewRuleStyle = lipgloss.NewStyle().
				Border(lipgloss.NormalBorder(), false, false, false, true).
				BorderForeground(lipgloss.Color("8")).
				PaddingLeft(2)
)

func (m model) View() string {
	if m.quitting {
		return ""
	}

	width := m.contentWidth()
	header := m.renderHeader(width)
	search := searchStyle.Width(width).Render(m.input.View())
	body := m.renderBody(width)
	footer := m.renderFooter(width)

	return appStyle.Render(strings.Join([]string{header, search, body, footer}, "\n"))
}

func (m model) contentWidth() int {
	width := m.width - pageHPadding*2
	if width < minWidth {
		return minWidth
	}
	return width
}

func (m model) bodyHeight() int {
	height := m.height - 8
	if height < listMinHeight {
		return listMinHeight
	}
	return height
}

func (m model) renderHeader(width int) string {
	left := titleStyle.Render("Ghostty Themes")

	current := m.currentTheme
	if current == "" {
		current = "none"
	}

	meta := fmt.Sprintf("%d themes  current: %s", len(m.themes), current)
	if m.color != "" && m.color != "all" {
		meta = fmt.Sprintf("%s  color: %s", meta, m.color)
	}
	right := headerMetaStyle.Render(truncate(meta, width-lipgloss.Width(left)-1))

	gap := width - lipgloss.Width(left) - lipgloss.Width(right)
	if gap < 1 {
		gap = 1
	}

	return left + strings.Repeat(" ", gap) + right
}

func (m model) renderBody(width int) string {
	height := m.bodyHeight()
	if width < previewAt {
		return m.renderList(width, height)
	}

	listWidth := width - previewWidth - 2
	if listWidth < minWidth {
		return m.renderList(width, height)
	}

	list := m.renderList(listWidth, height)
	preview := previewRuleStyle.
		Width(previewWidth).
		Height(height).
		Render(m.renderPreview(previewWidth - 3))

	return lipgloss.JoinHorizontal(lipgloss.Top, list, "  ", preview)
}

func (m model) renderList(width, height int) string {
	if m.loading {
		return centerMessage(width, height, "Loading themes...")
	}

	if len(m.filtered) == 0 {
		query := strings.TrimSpace(m.input.Value())
		if query == "" {
			return centerMessage(width, height, "No themes found")
		}
		return centerMessage(width, height, fmt.Sprintf("No themes match %q", query))
	}

	start := 0
	if m.selected >= height {
		start = m.selected - height + 1
	}
	end := start + height
	if end > len(m.filtered) {
		end = len(m.filtered)
	}

	lines := make([]string, 0, height)
	for i := start; i < end; i++ {
		lines = append(lines, m.renderThemeRow(m.filtered[i], width, i == m.selected))
	}
	for len(lines) < height {
		lines = append(lines, "")
	}

	return lipgloss.NewStyle().Width(width).Height(height).Render(strings.Join(lines, "\n"))
}

func (m model) renderThemeRow(theme ghostty.Theme, width int, selected bool) string {
	prefix := "  "
	if selected {
		prefix = "> "
	}

	fixedWidth := len(prefix) + sourceColWidth + stateColWidth + 2
	nameWidth := width - fixedWidth
	if nameWidth < 8 {
		nameWidth = 8
	}

	source := mutedStyle.Width(sourceColWidth).Render(theme.Source)
	state := ""
	if strings.EqualFold(theme.Name, m.currentTheme) {
		state = "current"
	}
	state = currentStyle.Width(stateColWidth).Render(state)

	line := prefix +
		padRight(truncate(theme.Name, nameWidth), nameWidth) +
		" " + source +
		" " + state

	if selected {
		return selectedStyle.Width(width).Render(line)
	}
	return lipgloss.NewStyle().Width(width).Render(line)
}

func (m model) renderFooter(width int) string {
	status := m.footerMessage()
	if status == "" {
		status = "Ready"
	}

	statusStyle := mutedStyle.Width(width)
	if m.message != "" && !isInformationalMessage(m.message) {
		statusStyle = errorStyle.Width(width)
	}

	help := mutedStyle.Render("Enter save   Esc cancel   / search   Ctrl+R reload")
	return statusStyle.Render(truncate(status, width)) + "\n" + help
}

func isInformationalMessage(message string) bool {
	return strings.HasPrefix(message, "Loading") ||
		strings.HasPrefix(message, "Reloading") ||
		strings.HasPrefix(message, "Saving")
}

func (m model) renderPreview(width int) string {
	theme, ok := m.selectedTheme()
	if !ok {
		return mutedStyle.Render("No theme selected")
	}

	lines := []string{
		titleStyle.Render(truncate(theme.Name, width)),
		mutedStyle.Render(theme.Source),
		"",
		sampleLine(theme, width),
		"",
		labelValue("Foreground", theme.Colors.Foreground, width),
		labelValue("Background", theme.Colors.Background, width),
		labelValue("Cursor", theme.Colors.Cursor, width),
		"",
		mutedStyle.Render("Palette"),
		paletteLine(theme.Colors.Palette, 0, 8),
		paletteLine(theme.Colors.Palette, 8, 16),
	}

	return strings.Join(lines, "\n")
}

func sampleLine(theme ghostty.Theme, width int) string {
	fg := theme.Colors.Foreground
	bg := theme.Colors.Background
	if fg == "" {
		fg = "15"
	}
	if bg == "" {
		bg = "0"
	}

	return lipgloss.NewStyle().
		Foreground(lipgloss.Color(fg)).
		Background(lipgloss.Color(bg)).
		Padding(0, 1).
		Width(width).
		Render("AaBbCc 0123456789")
}

func labelValue(label, value string, width int) string {
	if value == "" {
		value = "not declared"
	}
	labelWidth := 11
	valueWidth := width - labelWidth - 1
	if valueWidth < 4 {
		valueWidth = 4
	}
	labelText := mutedStyle.Width(labelWidth).Render(label)
	return labelText + " " + truncate(value, valueWidth)
}

func paletteLine(palette map[int]string, start, end int) string {
	if len(palette) == 0 {
		return mutedStyle.Render("not declared")
	}

	indexes := make([]int, 0, end-start)
	for i := start; i < end; i++ {
		if _, ok := palette[i]; ok {
			indexes = append(indexes, i)
		}
	}
	sort.Ints(indexes)

	if len(indexes) == 0 {
		return mutedStyle.Render("not declared")
	}

	parts := make([]string, 0, len(indexes))
	for _, index := range indexes {
		parts = append(parts, lipgloss.NewStyle().
			Background(lipgloss.Color(palette[index])).
			Render("  "))
	}
	return strings.Join(parts, " ")
}

func centerMessage(width, height int, message string) string {
	if height < 1 {
		height = 1
	}
	top := height / 2
	lines := make([]string, 0, height)
	for i := 0; i < top; i++ {
		lines = append(lines, "")
	}
	lines = append(lines, mutedStyle.Width(width).Align(lipgloss.Center).Render(truncate(message, width)))
	for len(lines) < height {
		lines = append(lines, "")
	}
	return lipgloss.NewStyle().Width(width).Height(height).Render(strings.Join(lines, "\n"))
}

func truncate(value string, limit int) string {
	if limit < 1 {
		return ""
	}
	runes := []rune(value)
	if len(runes) <= limit {
		return value
	}
	if limit <= 3 {
		return strings.Repeat(".", limit)
	}
	return string(runes[:limit-3]) + "..."
}

func padRight(value string, width int) string {
	current := lipgloss.Width(value)
	if current >= width {
		return truncate(value, width)
	}
	return value + strings.Repeat(" ", width-current)
}
