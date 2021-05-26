package render

type InternationalizationString struct {
	Default string
	Map     map[string]string
}

func (s InternationalizationString) String(languages []Language) string {
	for _, language := range languages {
		if str, ok := s.Map[language.Language]; ok {
			return str
		}
	}
	return s.Default
}
