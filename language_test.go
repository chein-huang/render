package render_test

import (
	"testing"

	"github.com/chein-huang/render"
	"github.com/stretchr/testify/require"
)

func TestGetLanguages(t *testing.T) {
	assert := require.New(t)
	languages := render.GetLanguages("en;q=0.9, zh, fr;q= 0.8")
	assert.Equal(
		[]render.Language{
			{
				"zh",
				1,
			}, {
				"en",
				0.9,
			}, {
				"fr",
				0.8,
			},
		},
		languages,
	)

	languages = render.GetLanguages("en;q=se")
	assert.Equal(
		[]render.Language{},
		languages,
	)

	languages = render.GetLanguages("")
	assert.Equal(
		[]render.Language{},
		languages,
	)
}
