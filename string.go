package render

type I18nResource struct {
	Desc    string            `json:"desc"`
	Default string            `json:"default"`
	Map     map[string]string `json:"map"`
}

func (s I18nResource) String(languages []Language) string {
	for _, language := range languages {
		if str, ok := s.Map[language.Language]; ok {
			return str
		}
	}
	return s.Default
}

type I18nID string

func (i I18nID) Resource() *I18nResource {
	return MustResource(string(i))
}

func (i I18nID) String(ls []Language) string {
	return MustResource(string(i)).String(ls)
}
