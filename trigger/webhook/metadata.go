package webhook

import "github.com/postqode/qodeflow-core/data/metadata"

// Settings
type Settings struct {
	Port int `md:"port"`
}

// HandlerSettings
type HandlerSettings struct {
	Method string `md:"method"`
	Path   string `md:"path"`
}

// Output
type Output struct {
	Method      string            `md:"method"`
	PathParams  map[string]string `md:"pathParams"`
	QueryParams map[string]string `md:"queryParams"`
	Headers     map[string]string `md:"headers"`
	Content     interface{}       `md:"content"`
}

// Reply
type Reply struct {
	Code int         `md:"code"`
	Data interface{} `md:"data"`
}

func (r *Reply) FromMap(values map[string]interface{}) error {
	return metadata.MapToStruct(values, r, true)
}
