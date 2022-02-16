package render

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

var (
	globalMap = map[string]*I18nResource{
		string(UnknownErrResponseMsg): {
			Default: GetTemplate("UnknownErrResponseMsg", "系统错误，请稍后重试！"),
			Map: map[string]*template.Template{
				"zh_CN": GetTemplate("UnknownErrResponseMsg.zh_CN", "系统错误，请稍后重试！"),
				"en_US": GetTemplate("UnknownErrResponseMsg.en_US", "System error, please try again later!"),
			},
		},
	}
)

type I18nStringResource struct {
	Desc    string            `json:"desc"`
	Default string            `json:"default"`
	Map     map[string]string `json:"map"`
}

func Resource(id string) (*I18nResource, bool) {
	s, ok := globalMap[id]
	return s, ok
}

func MustResource(id string) *I18nResource {
	s, ok := globalMap[id]
	if !ok {
		return nil
	}
	return s
}

func String(id string, ls []Language, args map[string]interface{}) (string, bool) {
	s, ok := Resource(id)
	if !ok {
		return "", ok
	}

	return s.String(ls, args), true
}

func MustString(id string, ls []Language, args map[string]interface{}) string {
	s := MustResource(id)
	return s.String(ls, args)
}

func SetResource(prefix, id string, resource *I18nStringResource) {
	key := prefix + "." + id
	if _, ok := globalMap[key]; !ok {
		newResource := &I18nResource{
			Desc:    resource.Desc,
			Default: GetTemplate(key, resource.Default),
			Map:     map[string]*template.Template{},
		}
		for mk, mv := range resource.Map {
			newResource.Map[mk] = GetTemplate(key+"."+mk, mv)
		}

		globalMap[key] = newResource
	} else {
		panic(fmt.Sprintf("key: %v is already in map", key))
	}
}

func SetResources(prefix string, resources map[string]*I18nStringResource) {
	for resID, resource := range resources {
		key := prefix + "." + resID
		if _, ok := globalMap[key]; !ok {
			newResource := &I18nResource{
				Desc:    resource.Desc,
				Default: GetTemplate(key, resource.Default),
				Map:     map[string]*template.Template{},
			}
			for mk, mv := range resource.Map {
				newResource.Map[mk] = GetTemplate(key+"."+mk, mv)
			}

			globalMap[key] = newResource
		} else {
			panic(fmt.Sprintf("key: %v is already in map", key))
		}
	}
}

func SetByJsonDecoder(prefix string, d *json.Decoder) error {
	m := map[string]*I18nStringResource{}
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

func GetTemplate(name, text string) *template.Template {
	return template.Must(template.New(name).Parse(text))
}
