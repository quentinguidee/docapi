package types

import (
	"fmt"
	"strings"
)

type RefType string

var (
	RefSchema   RefType = "schemas"
	RefResponse RefType = "responses"
)

type (
	Format struct {
		Openapi    string                  `json:"openapi" yaml:"openapi"`
		Info       FormatInfo              `json:"info" yaml:"info"`
		Servers    []FormatServer          `json:"servers,omitempty" yaml:"servers,omitempty"`
		Paths      map[string]FormatRoutes `json:"paths,omitempty" yaml:"paths,omitempty"`
		Components FormatComponents        `json:"components,omitempty" yaml:"components,omitempty"`
	}

	FormatInfo struct {
		Title       string `json:"title" yaml:"title"`
		Description string `json:"description" yaml:"description"`
		Version     string `json:"version" yaml:"version"`
	}

	FormatServer struct {
		Url         string                          `json:"url" yaml:"url"`
		Description string                          `json:"description,omitempty" yaml:"description,omitempty"`
		Variables   map[string]FormatServerVariable `json:"variables,omitempty" yaml:"variables,omitempty"`
	}

	FormatServerVariable struct {
		Default     string `json:"default" yaml:"default"`
		Description string `json:"description,omitempty" yaml:"description,omitempty"`
	}

	FormatRoutes map[string]FormatRoute

	FormatRoute struct {
		OperationId string                    `json:"operationId,omitempty" yaml:"operationId,omitempty"`
		Summary     string                    `json:"summary,omitempty" yaml:"summary,omitempty"`
		Tags        []string                  `json:"tags,omitempty" yaml:"tags,omitempty"`
		Description string                    `json:"description,omitempty" yaml:"description,omitempty"`
		Parameters  []FormatParameter         `json:"parameters,omitempty" yaml:"parameters,omitempty"`
		RequestBody FormatRequestBody         `json:"requestBody,omitempty" yaml:"requestBody,omitempty"`
		Responses   map[string]FormatResponse `json:"responses" yaml:"responses"`
	}

	FormatRequestBody struct {
		Description string                   `json:"description,omitempty" yaml:"description,omitempty"`
		Required    bool                     `json:"required,omitempty" yaml:"required,omitempty"`
		Content     map[string]FormatContent `json:"content,omitempty" yaml:"content,omitempty"`
	}

	FormatParameter struct {
		In          string       `json:"in,omitempty" yaml:"in,omitempty"`
		Name        string       `json:"name,omitempty" yaml:"name,omitempty"`
		Description string       `json:"description,omitempty" yaml:"description,omitempty"`
		Required    bool         `json:"required,omitempty" yaml:"required,omitempty"`
		Schema      FormatSchema `json:"schema,omitempty" yaml:"schema,omitempty"`
	}

	FormatResponse struct {
		Ref         Ref                      `json:"$ref,omitempty" yaml:"$ref,omitempty"`
		Description string                   `json:"description,omitempty" yaml:"description,omitempty"`
		Content     map[string]FormatContent `json:"content,omitempty" yaml:"content,omitempty"`
	}

	FormatContent struct {
		Schema FormatSchema `json:"schema,omitempty" yaml:"schema,omitempty"`
	}

	FormatSchema struct {
		Type       string                  `json:"type,omitempty" yaml:"type,omitempty"`
		Items      FormatItems             `json:"items,omitempty" yaml:"items,omitempty"`
		Properties map[string]FormatSchema `json:"properties,omitempty" yaml:"properties,omitempty"`
		Ref        Ref                     `json:"$ref,omitempty" yaml:"$ref,omitempty"`
	}

	FormatItems struct {
		Ref Ref `json:"$ref,omitempty" yaml:"$ref,omitempty"`
	}

	FormatComponents struct {
		Responses map[string]FormatResponse `json:"responses,omitempty" yaml:"responses,omitempty"`
		Schemas   map[string]FormatSchema   `json:"schemas,omitempty" yaml:"schemas,omitempty"`
	}

	Ref struct {
		Ref string `json:"$ref,omitempty" yaml:"$ref,omitempty"`
	}
)

func (f *Format) AddServer(server FormatServer) {
	f.Servers = append(f.Servers, server)
}

func (f *FormatServer) SetVariable(name string, variable FormatServerVariable) {
	if f.Variables == nil {
		f.Variables = map[string]FormatServerVariable{}
	}
	f.Variables[name] = variable
}

func (f *FormatRoute) SetResponse(code string, resp FormatResponse) {
	if f.Responses == nil {
		f.Responses = map[string]FormatResponse{}
	}
	f.Responses[code] = resp
}

func (f *FormatRoute) AddParameter(param FormatParameter) {
	f.Parameters = append(f.Parameters, param)
}

func (f *FormatSchema) SetProperty(name string, schema FormatSchema) {
	if f.Properties == nil {
		f.Properties = map[string]FormatSchema{}
	}
	f.Properties[name] = schema
}

func (f *FormatComponents) SetResponse(code string, resp FormatResponse) {
	if f.Responses == nil {
		f.Responses = map[string]FormatResponse{}
	}
	f.Responses[code] = resp
}

func (f *FormatComponents) SetSchema(name string, schema FormatSchema) {
	if f.Schemas == nil {
		f.Schemas = map[string]FormatSchema{}
	}
	f.Schemas[name] = schema
}

func CreateRef(tp RefType, name string) Ref {
	return Ref{
		Ref: fmt.Sprintf("#/components/%s/%s", tp, name),
	}
}

func (f Ref) Name() string {
	if f.Ref == "" {
		return ""
	}
	r := strings.Split(f.Ref, "/")
	return r[len(r)-1]
}

func (f Ref) MarshalJSON() ([]byte, error) {
	return []byte(f.Ref), nil
}

func (f Ref) MarshalYAML() (interface{}, error) {
	return f.Ref, nil
}

func (f *Format) GetReferencedComponents() []string {
	var schemas []string
	for _, route := range f.Paths {
		for _, resp := range route {
			schemas = append(schemas, resp.GetReferencedComponents()...)
		}
	}
	schemas = append(schemas, f.Components.GetReferencedComponents()...)
	return schemas
}

func (f *FormatRoute) GetReferencedComponents() []string {
	var schemas []string
	for _, param := range f.Parameters {
		schemas = append(schemas, param.GetReferencedComponents()...)
	}
	for _, resp := range f.Responses {
		schemas = append(schemas, resp.GetReferencedComponents()...)
	}
	if f.RequestBody.Content != nil {
		for _, content := range f.RequestBody.Content {
			schemas = append(schemas, content.GetReferencedComponents()...)
		}
	}
	return schemas
}

func (f *FormatParameter) GetReferencedComponents() []string {
	return f.Schema.GetReferencedComponents()
}

func (f *FormatResponse) GetReferencedComponents() []string {
	if f.Ref.Name() != "" {
		return []string{f.Ref.Name()}
	}

	var schemas []string
	for _, content := range f.Content {
		schemas = append(schemas, content.GetReferencedComponents()...)
	}
	return schemas
}

func (f *FormatContent) GetReferencedComponents() []string {
	return f.Schema.GetReferencedComponents()
}

func (f *FormatSchema) GetReferencedComponents() []string {
	if f.Ref.Name() != "" {
		return []string{f.Ref.Name()}
	}

	var schemas []string
	if f.Items.Ref.Name() != "" {
		schemas = append(schemas, f.Items.Ref.Name())
	}
	for _, schema := range f.Properties {
		schemas = append(schemas, schema.GetReferencedComponents()...)
	}
	return schemas
}

func (f *FormatComponents) GetReferencedComponents() []string {
	var schemas []string
	for _, resp := range f.Responses {
		schemas = append(schemas, resp.GetReferencedComponents()...)
	}
	for _, schema := range f.Schemas {
		schemas = append(schemas, schema.GetReferencedComponents()...)
	}
	return schemas
}
