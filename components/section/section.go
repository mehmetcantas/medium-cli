package section

import (
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/mehmetcantas/medium-cli/components/constants"
	"github.com/mehmetcantas/medium-cli/components/table"
	"github.com/mehmetcantas/medium-cli/config"
	"github.com/mehmetcantas/medium-cli/ui/screencontext"
)

type Model struct {
	Id        int
	Config    config.SectionConfig
	Ctx       *screencontext.ScreenContext
	Spinner   spinner.Model
	IsLoading bool
	Table     table.Model
	Type      string
}

type Section interface {
	Id() int
	Update(msg tea.Msg) (Section, tea.Cmd)
	View() string
	NumRows() int
	GetCurrRow() interface{}
	NextRow() int
	PrevRow() int
	FetchSectionRows() tea.Cmd
	GetIsLoading() bool
	GetSectionColumns() []table.Column
	BuildRows() []table.Row
	UpdateScreenContext(ctx *screencontext.ScreenContext)
}

func (m *Model) CreateNextTickCmd(nextTickCmd tea.Cmd) tea.Cmd {
	if m == nil || nextTickCmd == nil {
		return nil
	}
	return func() tea.Msg {
		return SectionTickMsg{
			SectionId:       m.Id,
			InternalTickMsg: nextTickCmd(),
			Type:            m.Type,
		}
	}
}

func (m *Model) GetDimensions() constants.Dimensions {
	return constants.Dimensions{
		Width: m.Ctx.MainContentWidth - lipgloss.NewStyle().
			Padding(0, 1).GetHorizontalPadding(),
		Height: m.Ctx.MainContentHeight - 2,
	}
}

func (m *Model) UpdateScreenContext(ctx *screencontext.ScreenContext) {
	oldDimensions := m.GetDimensions()
	m.Ctx = ctx
	newDimensions := m.GetDimensions()
	m.Table.SetDimensions(newDimensions)

	if oldDimensions.Height != newDimensions.Height || oldDimensions.Width != newDimensions.Width {
		m.Table.SyncViewPortContent()
	}
}

type SectionMsg interface {
	GetSectionId() int
	GetSectionType() string
}

type SectionRowsFetchedMsg struct {
	SectionId   int
	Placeholder []interface{}
}

func (msg *SectionRowsFetchedMsg) GetSectionId() int {
	return msg.SectionId
}

type SectionTickMsg struct {
	SectionId       int
	InternalTickMsg tea.Msg
	Type            string
}

func (msg SectionTickMsg) GetSectionType() string {
	return msg.Type
}
func (m *Model) NextRow() int {
	if m != nil && len(m.Table.Rows) > 1 {
		return m.Table.NextItem()
	}
	return 0
}
func (m *Model) PrevRow() int {
	if m != nil && len(m.Table.Rows) > 1 {
		return m.Table.PrevItem()
	}
	return 0
}

func (m *Model) GetIsLoading() bool {
	return m.IsLoading
}
