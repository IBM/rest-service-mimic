package handlers

import (
	"bytes"
	"html/template"
	"reflect"

	"github.com/Masterminds/sprig"
	"golang.org/x/crypto/bcrypt"
)

type ResponsePayloadGenerator struct {
	Request map[string]interface{}
	Cache   map[string]interface{}
}

func (generator ResponsePayloadGenerator) Generate(data interface{}) interface{} {
	if reflect.ValueOf(data).Kind() == reflect.Slice {
		return generator.handleSlice(reflect.ValueOf(data))
	} else if reflect.ValueOf(data).Kind() == reflect.Map {
		return generator.handleMap(reflect.ValueOf(data))
	} else if reflect.ValueOf(data).Kind() == reflect.String {
		var tpl bytes.Buffer
		fieldTemplate := template.Must(template.New("field").Funcs(sprig.FuncMap()).Funcs(map[string]interface{}{"hashFromPassword": hashFromPassword}).Parse(reflect.ValueOf(data).String()))

		fieldTemplate.Execute(&tpl, generator)
		return tpl.String()
	}

	return data
}

func (generator ResponsePayloadGenerator) handleSlice(data reflect.Value) []interface{} {
	returnSlice := make([]interface{}, data.Len())
	for i := 0; i < data.Len(); i++ {
		returnSlice[i] = generator.Generate(data.Index(i).Interface())
	}

	return returnSlice
}

func (generator ResponsePayloadGenerator) handleMap(data reflect.Value) map[string]interface{} {
	returnMap := make(map[string]interface{})
	for _, key := range data.MapKeys() {
		returnMap[key.String()] = generator.Generate(data.MapIndex(key).Interface())
	}

	return returnMap
}

func hashFromPassword(raw string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(raw), 14)
	return string(bytes), err
}
