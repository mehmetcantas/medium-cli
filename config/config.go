package config

type ViewType string

const (
	PlaceholderView ViewType = "placeholder"
	OtherView       ViewType = "other"
)

type SectionConfig struct {
	Title   string
	Filters string
	Limit   *int `yaml:"limit,omitempty"`
}

type PreviewConfig struct {
	Open  bool
	Width int
}

type Defaults struct {
	Preview PreviewConfig `yaml:"preview"`
	View    ViewType      `yaml:"view"`
}

type Config struct {
	PlaceholderSections []SectionConfig `yaml:"placeholderSections"`
	OtherSections       []SectionConfig `yaml:"otherSections"`
	Defaults            Defaults        `yaml:"defaults"`
}

type configError struct {
	configDir string
	parser    ConfigParser
	err       error
}

type ConfigParser struct{}

func (p ConfigParser) getDefaultConfig() Config {
	return Config{
		Defaults: Defaults{
			Preview: PreviewConfig{
				Open:  true,
				Width: 50,
			},
			View: PlaceholderView,
		},
		PlaceholderSections: []SectionConfig{
			{
				Title: "Albums",
				// Buraya ön tanımlı filtrelerinizi yazabilirsiniz. Örneğin Github ile alakalı bir uygulama geliştiriyorsanız
				// Kendi yarattığınız issue'ları görüntülemek isteyebilirsiniz o zaman bu tab için aşağıdaki gibi bir ön tanımlı filtre oluşturabilirsiniz

				//is:open author:@me
				Filters: "albums",
			},
			{
				Title:   "Todos",
				Filters: "todos",
			},
		},
		OtherSections: []SectionConfig{
			{
				Title:   "Comments",
				Filters: "",
			},
			{
				Title:   "Photos",
				Filters: "",
			},
			{
				Title:   "Users",
				Filters: "",
			},
		},
	}
}
func initParser() ConfigParser {
	return ConfigParser{}
}
func ParseConfig() (Config, error) {
	parser := initParser()

	return parser.getDefaultConfig(), nil
}
