package handlers

import (
	"encoding/json"
	"reflect"
	"strconv"
)

type JsonEnvelope struct {
	Response interface{}       `json:"response"`
	Meta     map[string]string `json:"meta"`
	Error    interface{}       `json:"error"`
}

func (e *JsonEnvelope) ToJson() string {
	json, _ := json.Marshal(e.Response)
	return string(json)
}

func getCount(obj interface{}) int {
	var count int

	switch reflect.TypeOf(obj).Kind() {
	case reflect.Slice:
		count = reflect.ValueOf(obj).Len()
	default:
		count = 1
	}

	return count
}

func getMeta(objs interface{}) map[string]string {
	count := getCount(objs)
	return map[string]string{"count": strconv.Itoa(count)}
}

func NewJsonEnvelope(objs interface{}, meta map[string]string, error interface{}) JsonEnvelope {
	return JsonEnvelope{
		Response: objs,
		Meta:     meta,
		Error:    error,
	}
}

func NewJsonEnvelopeFromObj(objs interface{}) JsonEnvelope {
	return NewJsonEnvelope(objs, getMeta(objs), nil)
}

func NewJsonEnvelopeFromError(err error) JsonEnvelope {
	error := map[string]string{
		"message": err.Error(),
	}

	meta := map[string]string{"count": "1"}

	return NewJsonEnvelope("", meta, error)
}
