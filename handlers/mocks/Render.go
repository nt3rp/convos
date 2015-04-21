package mocks

import (
	"html/template"
	"net/http"

	"github.com/martini-contrib/render"
)

type Render struct {
	StatusCode int
	Response   interface{}
}

func (m *Render) JSON(status int, v interface{}) {
	m.StatusCode = status
	m.Response = v
}

func (m *Render) HTML(status int, name string, v interface{}, htmlOpt ...render.HTMLOptions) {

}

func (m *Render) XML(status int, v interface{}) {

}

func (m *Render) Data(status int, v []byte) {

}

func (m *Render) Error(status int) {

}

func (m *Render) Status(status int) {

}

func (m *Render) Redirect(location string, status ...int) {

}

func (m *Render) Template() *template.Template {
	return &template.Template{}
}

func (m *Render) Header() http.Header {
	return http.Header{}
}
