package universal

// API

type Localizations interface {
	Add(country string, language string, translation string) (Localization, error)
}
