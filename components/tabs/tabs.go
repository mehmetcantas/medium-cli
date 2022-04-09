package tabs

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/mehmetcantas/medium-cli/config"
	"github.com/mehmetcantas/medium-cli/ui/screencontext"
)

var (
	tabsBorderHeight  = 1
	tabsContentHeight = 2
	TabsHeight        = tabsBorderHeight + tabsContentHeight
	blue              = lipgloss.AdaptiveColor{Light: "#3498db", Dark: "#2980b9"}
	tab               = lipgloss.NewStyle().
				Faint(true).
				Padding(0, 2)
	activeTab = tab.
			Copy().
			Faint(false).
			Bold(true).
			Background(lipgloss.AdaptiveColor{Light: blue.Light, Dark: "#3498db"}).
			Foreground(lipgloss.AdaptiveColor{Light: "#2980b9", Dark: "#E2E1ED"})

	tabsRow = lipgloss.NewStyle().
		Height(tabsContentHeight).
		PaddingTop(1).
		PaddingBottom(0).
		BorderBottom(true).
		BorderStyle(lipgloss.ThickBorder()).
		BorderBottomForeground(lipgloss.AdaptiveColor{Light: blue.Light, Dark: blue.Dark})

	viewSwitcher = lipgloss.NewStyle()

	activeView = lipgloss.NewStyle().
			Foreground(lipgloss.AdaptiveColor{Light: "#242347", Dark: "#E2E1ED"}).
			MarginLeft(1).
			Bold(true).
			Background(lipgloss.AdaptiveColor{Light: blue.Light, Dark: "#39386b"})

	inactiveView = lipgloss.NewStyle().
			MarginLeft(1).
			Background(lipgloss.AdaptiveColor{Light: "#D9DCCF", Dark: "#2b2b40"}).
			Foreground(lipgloss.AdaptiveColor{Light: blue.Light, Dark: "#666CA6"})
)

type Model struct {
	CurrSectionId int
}

func NewModel() Model {
	return Model{
		CurrSectionId: 0,
	}
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	return m, nil
}

func (m Model) View(ctx screencontext.ScreenContext) string {
	sectionsConfigs := ctx.GetViewSectionsConfig()
	sectionTitles := make([]string, 0, len(sectionsConfigs))
	for _, section := range sectionsConfigs {
		sectionTitles = append(sectionTitles, section.Title)
	}

	var tabs []string
	for i, sectionTitle := range sectionTitles {
		if m.CurrSectionId == i {
			tabs = append(tabs, activeTab.Render(sectionTitle))
		} else {
			tabs = append(tabs, tab.Render(sectionTitle))
		}
	}

	viewSwitcher := m.renderViewSwitcher(ctx)
	tabsWidth := ctx.ScreenWidth - lipgloss.Width(viewSwitcher)
	renderedTabs := lipgloss.NewStyle().
		Width(tabsWidth).
		MaxWidth(tabsWidth).
		Render(lipgloss.JoinHorizontal(lipgloss.Top, tabs...))

	return tabsRow.Copy().
		Width(ctx.ScreenWidth).
		MaxWidth(ctx.ScreenWidth).
		Render(lipgloss.JoinHorizontal(lipgloss.Top, renderedTabs, viewSwitcher))
}
func (m *Model) SetCurrSectionId(id int) {
	m.CurrSectionId = id
}

func (m *Model) renderViewSwitcher(ctx screencontext.ScreenContext) string {
	var placeholderStyle lipgloss.Style //,otherStyle
	if ctx.View == config.PlaceholderView {
		placeholderStyle = activeView
		//otherStyle = inactiveView
	} else {
		placeholderStyle = inactiveView
		//otherStyle = activeView
	}

	placeholder := placeholderStyle.Render("[療Placeholder]")
	//other := otherStyle.Render("[ﭦ Other]")
	return viewSwitcher.Copy().
		Render(lipgloss.JoinHorizontal(lipgloss.Top, placeholder)) //, other
}
