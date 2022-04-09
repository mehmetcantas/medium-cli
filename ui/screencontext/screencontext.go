package screencontext

import "github.com/mehmetcantas/medium-cli/config"

type ScreenContext struct {
	ScreenHeight      int
	ScreenWidth       int
	MainContentWidth  int
	MainContentHeight int
	Config            *config.Config
	View              config.ViewType
}

func (ctx *ScreenContext) GetViewSectionsConfig() []config.SectionConfig {

	return ctx.Config.PlaceholderSections

}
