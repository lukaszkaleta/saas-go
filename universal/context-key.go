package universal

import "context"

type contextKey string

const languageKey contextKey = "language"
const currentUserKey contextKey = "current-user-id"

func CurrentUserId(ctx context.Context) *int64 {
	return ctx.Value(currentUserKey).(*int64)
}

func WithLanguage(ctx context.Context, lang string) context.Context {
	return context.WithValue(ctx, languageKey, lang)
}

func GetLanguage(ctx context.Context) string {
	if lang, ok := ctx.Value(languageKey).(string); ok {
		return lang
	}
	return "no" // default fallback
}
