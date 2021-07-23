package render

import (
	"strings"
	"text/template"
)

type I18nResource struct {
	Desc    string                        `json:"desc"`
	Default *template.Template            `json:"default"`
	Map     map[string]*template.Template `json:"map"`
}

func (s I18nResource) String(languages []Language, args map[string]interface{}) string {
	var t *template.Template
	for _, language := range languages {
		if templ, ok := s.Map[language.Language]; ok {
			t = templ
		}
	}
	if t == nil {
		t = s.Default
	}

	return TemplateString(t, args)
}

type I18nID string

func (i I18nID) Resource() *I18nResource {
	return MustResource(string(i))
}

func (i I18nID) String(ls []Language, args map[string]interface{}) string {
	return MustResource(string(i)).String(ls, args)
}

func TemplateString(t *template.Template, args map[string]interface{}) string {
	strbuf := strings.Builder{}
	if err := t.Execute(&strbuf, args); err != nil {
		panic(err)
	}
	return strbuf.String()
}
