package placeholdersection

import (
	"fmt"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/mehmetcantas/medium-cli/components/constants"
	"github.com/mehmetcantas/medium-cli/components/section"
	"github.com/mehmetcantas/medium-cli/components/table"
	"github.com/mehmetcantas/medium-cli/config"
	"github.com/mehmetcantas/medium-cli/ui/screencontext"
)

const SectionType = "placeholder"

var (
	updatedAtCellWidth = lipgloss.Width("sequi sint nihil reprehenderit dolor beatae")
	ContainerPadding   = 1

	containerStyle = lipgloss.NewStyle().
			Padding(0, ContainerPadding)

	spinnerStyle = lipgloss.NewStyle().Padding(0, 1)

	emptyStateStyle = lipgloss.NewStyle().
			Faint(true).
			PaddingLeft(1).
			MarginBottom(1)
)

type Model struct {
	Placeholders      []PlaceholderModel
	section           section.Model
	err               error
	placeholderClient *PlaceholderClient
}

func NewModel(id int, ctx *screencontext.ScreenContext, config config.SectionConfig) Model {
	placeholderClient := NewPlaceholderClient("https://jsonplaceholder.typicode.com")
	m := Model{
		Placeholders: []PlaceholderModel{},
		section: section.Model{
			Id:        id,
			Config:    config,
			Ctx:       ctx,
			Spinner:   spinner.Model{Spinner: spinner.Moon},
			IsLoading: true,
			Type:      SectionType,
		},
		placeholderClient: placeholderClient,
		err:               nil,
	}

	m.section.Table = table.NewModel(
		m.getDimensions(),
		m.GetSectionColumns(),
		m.BuildRows(),
		"Placeholder",
		emptyStateStyle.Render("No data found"),
		m.section.Config.Title,
	)

	return m
}

func (m Model) Update(msg tea.Msg) (section.Section, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case SectionPlaceholdersFetchedMsg:
		m.Placeholders = msg.Placeholders
		m.section.IsLoading = false
		m.section.Table.SetRows(m.BuildRows())
		m.err = msg.Err
	case section.SectionTickMsg:
		if m.section.IsLoading == false {
			return &m, nil
		}

		var internalTickCmd tea.Cmd
		m.section.Spinner, internalTickCmd = m.section.Spinner.Update(msg.InternalTickMsg)
		cmd = m.section.CreateNextTickCmd(internalTickCmd)
	}

	return &m, cmd
}

func (m *Model) getDimensions() constants.Dimensions {
	return constants.Dimensions{
		Width:  m.section.Ctx.MainContentWidth - containerStyle.GetHorizontalPadding(),
		Height: m.section.Ctx.MainContentHeight - 2,
	}
}

func (m *Model) View() string {
	var spinnerText string
	if m.section.IsLoading {
		spinnerText = lipgloss.JoinHorizontal(lipgloss.Top, spinnerStyle.Copy().Render(m.section.Spinner.View()), "Fetching data...")
	}

	if m.err != nil {
		spinnerText = fmt.Sprintf("Error while fetching data : %v", m.err)
	}

	return containerStyle.Copy().Render(m.section.Table.View(spinnerText))
}

func (m *Model) UpdateScreenContext(ctx *screencontext.ScreenContext) {
	oldDimensions := m.getDimensions()
	m.section.Ctx = ctx
	newDimensions := m.getDimensions()
	m.section.Table.SetDimensions(newDimensions)

	if oldDimensions.Height != newDimensions.Height || oldDimensions.Width != newDimensions.Width {
		m.section.Table.SyncViewPortContent()
	}
}

func (m *Model) GetSectionColumns() []table.Column {
	return []table.Column{
		{
			Title: "ID",
			Width: &updatedAtCellWidth,
		},
		{
			Title: "Title",
			Width: &updatedAtCellWidth,
		},
		{
			Title: "User ID",
			Width: &updatedAtCellWidth,
		},
	}
}

func (m *Model) BuildRows() []table.Row {
	var rows []table.Row
	for _, currPlaceholders := range m.Placeholders {
		placeholdersModel := Placeholder{Data: currPlaceholders, Width: m.getDimensions().Width}
		rows = append(rows, placeholdersModel.ToTableRow())
	}

	return rows
}

func (m *Model) NumRows() int {
	return len(m.Placeholders)
}

type SectionPlaceholdersFetchedMsg struct {
	SectionId    int
	Placeholders []PlaceholderModel
	Err          error
}

func (msg SectionPlaceholdersFetchedMsg) GetSectionId() int {
	return msg.SectionId
}

func (msg SectionPlaceholdersFetchedMsg) GetSectionType() string {
	return SectionType
}

func (m *Model) GetCurrRow() interface{} {
	if len(m.Placeholders) == 0 {
		return nil
	}

	placeholder := m.Placeholders[m.section.Table.GetCurrItem()]
	return placeholder
}
func (m *Model) NextRow() int {
	return m.section.Table.NextItem()
}

func (m *Model) PrevRow() int {
	return m.section.Table.PrevItem()
}

func (m *Model) FetchSectionRows() tea.Cmd {
	m.err = nil
	if m == nil {
		return nil
	}
	m.Placeholders = nil
	m.section.Table.ResetCurrItem()
	m.section.Table.Rows = nil
	m.section.IsLoading = true
	var cmds []tea.Cmd
	cmds = append(cmds, m.section.CreateNextTickCmd(spinner.Tick))

	cmds = append(cmds, func() tea.Msg {
		fetchedData := m.placeholderClient.Get(m.section.Config.Filters)

		if len(fetchedData) <= 0 {
			return SectionPlaceholdersFetchedMsg{
				SectionId:    m.section.Id,
				Placeholders: []PlaceholderModel{},
				Err:          nil,
			}
		}

		return SectionPlaceholdersFetchedMsg{
			SectionId:    m.section.Id,
			Placeholders: fetchedData,
		}
	})

	return tea.Batch(cmds...)
}

func (m *Model) Id() int {
	return m.section.Id
}

func (m *Model) GetIsLoading() bool {
	return m.section.IsLoading
}

func FetchAllSections(ctx screencontext.ScreenContext) (sections []section.Section, fetchAllCmd tea.Cmd) {
	sectionConfigs := ctx.Config.PlaceholderSections
	fetchIssuesCmds := make([]tea.Cmd, 0, len(sectionConfigs))
	sections = make([]section.Section, 0, len(sectionConfigs))
	for i, sectionConfig := range sectionConfigs {
		sectionModel := NewModel(i, &ctx, sectionConfig)
		sections = append(sections, &sectionModel)
		fetchIssuesCmds = append(fetchIssuesCmds, sectionModel.FetchSectionRows())
	}
	return sections, tea.Batch(fetchIssuesCmds...)
}
