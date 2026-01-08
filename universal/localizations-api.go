package universal

import "context"

// API

type Localizations interface {
	Add(ctx context.Context, country string, language string, translation string) (Localization, error)
}
