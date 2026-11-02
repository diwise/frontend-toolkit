package locale

import (
	"fmt"
	"iter"
	"slices"
	"strings"

	"github.com/BurntSushi/toml"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
)

type LocaleBundle interface {
	For(acceptLanguage string) Localizer
	ForAll(languages []string) Localizer

	// GetAll returns an iterator that provides all translations for a given key, starting
	// with the default (fallback) translation and omitting any translation key that would
	// return the fallback result
	GetAll(id string) iter.Seq2[string, string]

	Languages() []string
}

type Localizer interface {
	Get(string) string
	GetWithData(string, map[string]any) string
}

type impl struct {
	bundle    *i18n.Bundle
	languages []string
}

type wrapper struct {
	localizer *i18n.Localizer
}

func NewLocalizer(assetPath string, languages ...string) LocaleBundle {

	if len(languages) == 0 {
		panic("at least one language must be specified")
	}

	defaultLanguage := language.Make(languages[0])

	b := &impl{
		bundle:    i18n.NewBundle(defaultLanguage),
		languages: make([]string, 0, len(languages)),
	}

	b.bundle.RegisterUnmarshalFunc("toml", toml.Unmarshal)

	for _, lang := range languages {
		b.bundle.MustLoadMessageFile(fmt.Sprintf("%s/l10n/%s.toml", assetPath, lang))
		b.languages = append(b.languages, lang)
	}

	return b
}

// For implements [LocaleBundle].
func (l *impl) For(acceptLanguage string) Localizer {
	return l.ForAll([]string{acceptLanguage})
}

// ForAll implements [LocaleBundle].
func (l *impl) ForAll(languages []string) Localizer {
	var localizer *i18n.Localizer

	switch len(languages) {
	case 0:
		panic("forall requires at least 1 language tag")
	case 1:
		localizer = i18n.NewLocalizer(l.bundle, languages[0])
	case 2:
		localizer = i18n.NewLocalizer(l.bundle, languages[0], languages[1])
	case 3:
		localizer = i18n.NewLocalizer(l.bundle, languages[0], languages[1], languages[2])
	default:
		localizer = i18n.NewLocalizer(l.bundle, languages...)
	}

	return &wrapper{localizer: localizer}
}

// GetAll implements [LocaleBundle].
func (l *impl) GetAll(id string) iter.Seq2[string, string] {
	return func(yield func(lang, str string) bool) {
		w := l.For(l.languages[0])
		defaultTranslation := w.Get(id)

		// Yield the default translation first
		if yield(l.languages[0], defaultTranslation) {
			for _, lang := range l.languages[1:] {
				w = l.For(lang)
				translated := w.Get(id)

				// Only yield translated strings that differ from the default
				if strings.Compare(defaultTranslation, translated) != 0 {
					if !yield(lang, translated) {
						return
					}
				}
			}
		}
	}
}

// Languages implements [LocaleBundle].
func (l *impl) Languages() []string {
	return slices.Clone(l.languages)
}

// Get implements [Localizer].
func (w *wrapper) Get(id string) string {
	cfg := &i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID:    id,
			Other: id,
			One:   id,
		},
	}

	str, err := w.localizer.Localize(cfg)
	if err != nil {
		return id
	}

	return str
}

// GetWithData implements [Localizer].
func (w *wrapper) GetWithData(id string, data map[string]any) string {
	cfg := &i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID:    id,
			Other: id,
		},
		TemplateData: data,
	}

	str, err := w.localizer.Localize(cfg)
	if err != nil {
		return id
	}

	return str
}
