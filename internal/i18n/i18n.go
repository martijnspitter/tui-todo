package i18n

import (
	"embed"
	"encoding/json"

	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
)

//go:embed translations/*.json
var translationFiles embed.FS

type TranslationService struct {
	bundle    *i18n.Bundle
	localizer *i18n.Localizer
	language  string
}

func NewTranslationService(defaultLang string) (*TranslationService, error) {
	bundle := i18n.NewBundle(language.MustParse(defaultLang))
	bundle.RegisterUnmarshalFunc("json", json.Unmarshal)

	entries, err := translationFiles.ReadDir("translations")
	if err != nil {
		return nil, err
	}

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		data, err := translationFiles.ReadFile("translations/" + entry.Name())
		if err != nil {
			return nil, err
		}

		if _, err := bundle.ParseMessageFileBytes(data, entry.Name()); err != nil {
			return nil, err
		}
	}

	localizer := i18n.NewLocalizer(bundle, defaultLang)

	return &TranslationService{
		bundle:    bundle,
		localizer: localizer,
		language:  defaultLang,
	}, nil
}

func (t *TranslationService) SetLanguage(lang string) {
	t.language = lang
	t.localizer = i18n.NewLocalizer(t.bundle, lang)
}

func (t *TranslationService) GetCurrentLanguage() string {
	return t.language
}

func (t *TranslationService) Translate(id string, templateData map[string]interface{}) string {
	message, err := t.localizer.Localize(&i18n.LocalizeConfig{
		MessageID:    id,
		TemplateData: templateData,
	})

	if err != nil {
		return id
	}

	return message
}

func (t *TranslationService) T(id string) string {
	return t.Translate(id, nil)
}

func (t *TranslationService) Tf(id string, data map[string]interface{}) string {
	return t.Translate(id, data)
}
