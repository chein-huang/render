package render_test

import (
	"testing"

	. "github.com/chein-huang/render"
	"github.com/stretchr/testify/require"
)

func TestI18n(t *testing.T) {
	assert := require.New(t)
	i := I18nResource{
		Desc:    "test",
		Default: GetTemplate("test", "default {{.msg}}"),
	}

	assert.Equal("default <no value>", i.String(nil, nil))
	assert.Equal("default msg", i.String(nil, map[string]interface{}{
		"msg": "msg",
	}))
}
