package render

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

var (
	resourceMap = map[string]*I18nResource{
		string(UnknownErrResponseMsg): {
			Default: "System error, please try again later!",
			Map: map[string]string{
				"zh_CN": "系统错误，请稍后重试！",
				"zh":    "系统错误，请稍后重试！",
				"en_US": "System error, please try again later!",
				"en":    "System error, please try again later!",
			},
		},
	}
)

func Resource(id string) (*I18nResource, bool) {
	s, ok := resourceMap[id]
	return s, ok
}

func MustResource(id string) *I18nResource {
	s, ok := resourceMap[id]
	if !ok {
		panic(fmt.Sprintf("%v not found", id))
	}
	return s
}

func String(id string, ls []Language) (string, bool) {
	s, ok := Resource(id)
	if !ok {
		return "", ok
	}

	return s.String(ls), true
}

func MustString(id string, ls []Language) string {
	s := MustResource(id)
	return s.String(ls)
}

func SetResource(prefix, id string, s *I18nResource) {
	key := prefix + "." + id
	if _, ok := resourceMap[key]; !ok {
		resourceMap[key] = s
	} else {
		panic(fmt.Sprintf("key: %v is already in map", key))
	}
}

func SetResources(prefix string, m map[string]*I18nResource) {
	for k, v := range m {
		key := prefix + "." + k
		if _, ok := resourceMap[key]; !ok {
			resourceMap[key] = v
		} else {
			panic(fmt.Sprintf("key: %v is already in map", key))
		}
	}
}

func SetByJsonDecoder(prefix string, d *json.Decoder) error {
	m := map[string]*I18nResource{}
	err := d.Decode(&m)
	if err != nil {
		return err
	}

	SetResources(prefix, m)
	return nil
}

func SetByFilename(prefix *string, filename string) error {
	f, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer f.Close()

	if prefix == nil {
		p := strings.TrimSuffix(filepath.Base(filename), filepath.Ext(filename))
		prefix = &p
	}

	return SetByJsonDecoder(*prefix, json.NewDecoder(f))
}
