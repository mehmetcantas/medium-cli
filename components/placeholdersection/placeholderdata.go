package placeholdersection

import (
	"github.com/charmbracelet/lipgloss"
	"github.com/mehmetcantas/medium-cli/components/table"
	"github.com/mehmetcantas/medium-cli/pkg"
)

type Placeholder struct {
	Data  PlaceholderModel
	Width int
}

type PlaceholderModel struct {
	UserId int    `json:"userId"`
	Id     int    `json:"id"`
	Title  string `json:"title"`
}

func (p *Placeholder) ToTableRow() table.Row {
	return table.Row{
		p.renderId(),
		p.renderTitle(),
		p.renderUserId(),
	}
}

func (p *Placeholder) renderId() string {
	return lipgloss.NewStyle().Render(pkg.CastIntToStr(p.Data.Id))
}

func (p *Placeholder) renderUserId() string {
	return lipgloss.NewStyle().Render(pkg.CastIntToStr(p.Data.UserId))
}
func (p *Placeholder) renderTitle() string {
	title := pkg.TruncateString(p.Data.Title, 18)
	return lipgloss.NewStyle().Render(title)
}
func (p *Placeholder) renderStatus() string {
	return lipgloss.NewStyle().Foreground(lipgloss.AdaptiveColor{Light: "#42A0FA", Dark: "#42A0FA"}).Render("")
}
