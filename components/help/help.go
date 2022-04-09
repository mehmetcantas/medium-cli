package help

import (
	bbHelp "github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/mehmetcantas/medium-cli/pkg"
	"github.com/mehmetcantas/medium-cli/ui/screencontext"
)

var (
	blue         = lipgloss.AdaptiveColor{Light: "#3498db", Dark: "#2980b9"}
	FooterHeight = 3

	helpTextStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#3498db"))
	helpStyle     = lipgloss.NewStyle().
			Height(FooterHeight - 1).
			BorderTop(true).
			BorderStyle(lipgloss.NormalBorder()).
			BorderForeground(lipgloss.Color("#3498db"))
)

type Model struct {
	help bbHelp.Model
}

func NewModel() Model {
	help := bbHelp.New()

	help.Styles = bbHelp.Styles{
		ShortDesc:      helpTextStyle.Copy(),
		FullDesc:       helpTextStyle.Copy(),
		ShortSeparator: helpTextStyle.Copy(),
		FullSeparator:  helpTextStyle.Copy(),
		FullKey:        helpTextStyle.Copy(),
		ShortKey:       helpTextStyle.Copy(),
		Ellipsis:       helpTextStyle.Copy(),
	}

	return Model{
		help: help,
	}
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if key.Matches(msg, pkg.Keys.Help) {
			m.help.ShowAll = !m.help.ShowAll
		}
	}

	return m, nil
}

func (m *Model) View(ctx screencontext.ScreenContext) string {
	return helpStyle.Copy().Width(ctx.ScreenWidth).Render(m.help.View(pkg.Keys))
}

func (m *Model) SetWidth(width int) {
	m.help.Width = width
}
