package ui

import (
	"strings"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/mehmetcantas/medium-cli/components/help"
	"github.com/mehmetcantas/medium-cli/components/placeholdersection"
	"github.com/mehmetcantas/medium-cli/components/section"
	"github.com/mehmetcantas/medium-cli/components/tabs"
	"github.com/mehmetcantas/medium-cli/config"
	"github.com/mehmetcantas/medium-cli/pkg"
	"github.com/mehmetcantas/medium-cli/ui/screencontext"
)

type Model struct {
	tabs          tabs.Model
	ctx           screencontext.ScreenContext
	keys          pkg.KeyMap
	placeholders  []section.Section
	err           error
	currSectionId int
	help          help.Model
}
type initMsg struct {
	Config config.Config
}

type errMsg struct {
	error
}

func (e errMsg) Error() string { return e.error.Error() }

func NewModel() Model {
	tabsModel := tabs.NewModel()
	return Model{
		keys:          pkg.Keys,
		currSectionId: 0,
		help:          help.NewModel(),
		tabs:          tabsModel,
	}
}
func initScreen() tea.Msg {
	settings, err := config.ParseConfig()
	if err != nil {
		return errMsg{err}
	}

	return initMsg{Config: settings}
}
func (m Model) Init() tea.Cmd {
	return tea.Batch(initScreen, tea.EnterAltScreen)
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		cmd         tea.Cmd
		sidebarCmd  tea.Cmd
		helpCmd     tea.Cmd
		cmds        []tea.Cmd
		currSection = m.getCurrSection()
	)

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keys.PrevSection):
			prevSection := m.getSectionAt(m.getPrevSectionId())
			if prevSection != nil {
				m.setCurrSectionId(prevSection.Id())
				m.onViewedRowChanged()
			}

		case key.Matches(msg, m.keys.NextSection):
			nextSectionId := m.getNextSectionId()
			nextSection := m.getSectionAt(nextSectionId)
			if nextSection != nil {
				m.setCurrSectionId(nextSection.Id())
				m.onViewedRowChanged()
			}
		case key.Matches(msg, m.keys.Up):
			currSection.PrevRow()
			m.onViewedRowChanged()

		case key.Matches(msg, m.keys.Down):
			currSection.NextRow()
			m.onViewedRowChanged()
		case key.Matches(msg, m.keys.Quit):
			cmd = tea.Quit

		case key.Matches(msg, m.keys.SwitchView):
			m.ctx.View = m.switchSelectedView()
			m.syncMainContentWidth()
			m.setCurrSectionId(0)

			currSections := m.getCurrentViewSections()
			if len(currSections) == 0 {
				newSections, fetchSectionsCmds := m.fetchAllViewSections()
				m.setCurrentViewSections(newSections)
				cmd = fetchSectionsCmds
			}
			m.onViewedRowChanged()
		case key.Matches(msg, m.keys.Refresh):
			cmd = currSection.FetchSectionRows()

		}
	case initMsg:
		m.ctx.Config = &msg.Config
		m.ctx.View = m.ctx.Config.Defaults.View
		m.syncMainContentWidth()
		newSections, fetchSectionsCmds := m.fetchAllViewSections()
		m.setCurrentViewSections(newSections)
		cmd = fetchSectionsCmds
	case section.SectionMsg:
		cmd = m.updateRelevantSection(msg)

		if msg.GetSectionId() == m.currSectionId {
			switch msg.GetSectionType() {
			case placeholdersection.SectionType:
				m.onViewedRowChanged()
			}
		}
	case tea.WindowSizeMsg:
		m.onWindowSizeChanged(msg)

	case errMsg:
		m.err = msg
	}

	m.syncProgramContext()
	m.help, helpCmd = m.help.Update(msg)
	cmds = append(cmds, cmd, sidebarCmd, helpCmd)
	return &m, tea.Batch(cmds...)
}

func (m Model) View() string {
	if m.err != nil {
		return m.err.Error()
	}

	if m.ctx.Config == nil {
		return "Reading config...\n"
	}

	s := strings.Builder{}
	s.WriteString(m.tabs.View(m.ctx))
	s.WriteString("\n")
	currSection := m.getCurrSection()
	mainContent := ""
	if currSection != nil {
		mainContent = lipgloss.JoinHorizontal(
			lipgloss.Top,
			m.getCurrSection().View(),
		)
	} else {
		mainContent = "No data found"
	}
	s.WriteString(mainContent)
	s.WriteString("\n")
	s.WriteString(m.help.View(m.ctx))
	return s.String()
}

func (m *Model) setCurrSectionId(newSectionId int) {
	m.currSectionId = newSectionId
	m.tabs.SetCurrSectionId(newSectionId)
}

func (m *Model) onViewedRowChanged() {

}
func (m *Model) getSectionAt(id int) section.Section {
	sections := m.getCurrentViewSections()
	if len(sections) <= id {
		return nil
	}
	return sections[id]
}
func (m *Model) onWindowSizeChanged(msg tea.WindowSizeMsg) {
	m.help.SetWidth(msg.Width)
	m.ctx.ScreenWidth = msg.Width
	m.ctx.ScreenHeight = msg.Height
	m.ctx.MainContentHeight = msg.Height - tabs.TabsHeight - help.FooterHeight
	m.syncMainContentWidth()
}

func (m *Model) syncProgramContext() {
	for _, section := range m.getCurrentViewSections() {
		section.UpdateScreenContext(&m.ctx)
	}
}
func (m *Model) syncMainContentWidth() {
	m.ctx.MainContentWidth = m.ctx.ScreenWidth
}

func (m *Model) getCurrSection() section.Section {
	sections := m.getCurrentViewSections()
	if len(sections) == 0 {
		return nil
	}
	return sections[m.currSectionId]
}

func (m *Model) getPrevSectionId() int {
	sectionsConfigs := m.ctx.GetViewSectionsConfig()
	m.currSectionId = (m.currSectionId - 1) % len(sectionsConfigs)
	if m.currSectionId < 0 {
		m.currSectionId += len(sectionsConfigs)
	}

	return m.currSectionId
}

func (m *Model) getNextSectionId() int {
	return (m.currSectionId + 1) % len(m.ctx.GetViewSectionsConfig())
}

func (m *Model) getCurrentViewSections() []section.Section {
	if m.ctx.View == config.PlaceholderView {
		return m.placeholders
	} else {
		return []section.Section{}
	}
}
func (m *Model) fetchAllViewSections() ([]section.Section, tea.Cmd) {
	if m.ctx.View == config.PlaceholderView {
		return placeholdersection.FetchAllSections(m.ctx)
	} else {
		return []section.Section{}, nil
	}

}
func (m *Model) setCurrentViewSections(newSections []section.Section) {
	if m.ctx.View == config.PlaceholderView {
		m.placeholders = newSections
	}
}

func (m *Model) updateRelevantSection(msg section.SectionMsg) (cmd tea.Cmd) {
	var updatedSection section.Section

	switch msg.GetSectionType() {
	case placeholdersection.SectionType:
		updatedSection, cmd = m.placeholders[msg.GetSectionId()].Update(msg)
		m.placeholders[msg.GetSectionId()] = updatedSection
	}

	return cmd
}
func (m *Model) switchSelectedView() config.ViewType {
	if m.ctx.View == config.PlaceholderView {
		return config.PlaceholderView
	} else {
		return config.OtherView
	}
}
