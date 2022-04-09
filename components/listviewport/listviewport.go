package listviewport

import (
	"fmt"

	"github.com/charmbracelet/bubbles/viewport"
	"github.com/charmbracelet/lipgloss"
	"github.com/mehmetcantas/medium-cli/components/constants"
	"github.com/mehmetcantas/medium-cli/pkg"
)

var (
	pagerHeight = 2
	blue        = lipgloss.AdaptiveColor{Light: "#3498db", Dark: "#2980b9"}
	pagerStyle  = lipgloss.NewStyle().
			Height(pagerHeight).
			MaxHeight(pagerHeight).
			PaddingTop(1).
			Bold(true).
			Foreground(lipgloss.AdaptiveColor{Light: blue.Light, Dark: "#3498db"})
)

type Model struct {
	viewport       viewport.Model
	topBoundId     int
	bottomBoundId  int
	currId         int
	ListItemHeight int
	NumItems       int
	TabName        string
	ItemTypeLabel  string
}

func NewModel(dimensions constants.Dimensions, itemTypeLabel string, numItems, listItemHeight int, tabName string) Model {
	model := Model{
		NumItems:       numItems,
		ListItemHeight: listItemHeight,
		currId:         0,
		viewport: viewport.Model{
			Width:  dimensions.Width,
			Height: dimensions.Height - pagerHeight,
		},
		topBoundId:    0,
		ItemTypeLabel: itemTypeLabel,
		TabName:       tabName,
	}
	model.bottomBoundId = pkg.Min(model.NumItems-1, model.getNumPrsPerPage()-1)
	return model
}

func (m *Model) SetNumItems(numItems int) {
	m.NumItems = numItems
	m.bottomBoundId = pkg.Min(m.NumItems-1, m.getNumPrsPerPage()-1)
}

func (m *Model) SyncViewPort(content string) {
	m.viewport.SetContent(content)
}

func (m *Model) getNumPrsPerPage() int {
	return m.viewport.Height / m.ListItemHeight
}

func (m *Model) ResetCurrItem() {
	m.currId = 0
}

func (m *Model) GetCurrItem() int {
	return m.currId
}

func (m *Model) NextItem() int {
	atBottomOfViewport := m.currId >= m.bottomBoundId
	if atBottomOfViewport {
		m.topBoundId += 1
		m.bottomBoundId += 1
		m.viewport.LineDown(m.ListItemHeight)
	}

	newId := pkg.Min(m.currId+1, m.NumItems-1)
	newId = pkg.Max(newId, 0)
	m.currId = newId
	return m.currId
}

func (m *Model) PrevItem() int {
	atTopOfViewport := m.currId < m.topBoundId
	if atTopOfViewport {
		m.topBoundId -= 1
		m.bottomBoundId -= 1
		m.viewport.LineUp(m.ListItemHeight)
	}

	m.currId = pkg.Max(m.currId-1, 0)
	return m.currId
}

func (m *Model) SetDimensions(dimensions constants.Dimensions) {
	m.viewport.Height = dimensions.Height - pagerHeight
	m.viewport.Width = dimensions.Width
}

func (m *Model) View() string {
	pagerContent := ""
	if m.NumItems > 0 {
		pagerContent = fmt.Sprintf(
			"%s %v/%v",
			m.TabName,
			m.currId+1,
			m.NumItems,
		)
	}
	viewport := m.viewport.View()
	pager := pagerStyle.Copy().Render(pagerContent)
	return lipgloss.NewStyle().
		Width(m.viewport.Width).
		MaxWidth(m.viewport.Width).
		Render(lipgloss.JoinVertical(
			lipgloss.Left,
			viewport,
			pager,
		))
}
