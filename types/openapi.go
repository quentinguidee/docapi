package types

type (
	Format struct {
		Openapi    string                  `json:"openapi" yaml:"openapi"`
		Info       FormatInfo              `json:"info" yaml:"info"`
		Servers    []FormatServer          `json:"servers,omitempty" yaml:"servers,omitempty"`
		Paths      map[string]FormatRoutes `json:"paths,omitempty" yaml:"paths,omitempty"`
		Components FormatComponents        `json:"components,omitempty" yaml:"components,omitempty"`
	}

	FormatInfo struct {
		Title   string `json:"title" yaml:"title"`
		Desc    string `json:"description" yaml:"description"`
		Version string `json:"version" yaml:"version"`
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
		Summary     string                    `json:"summary,omitempty" yaml:"summary,omitempty"`
		Tags        []string                  `json:"tags,omitempty" yaml:"tags,omitempty"`
		Description string                    `json:"description,omitempty" yaml:"description,omitempty"`
		Servers     []FormatServer            `json:"servers,omitempty" yaml:"servers,omitempty"`
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
		Ref         string                   `json:"$ref,omitempty" yaml:"$ref,omitempty"`
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
		Ref        string                  `json:"$ref,omitempty" yaml:"$ref,omitempty"`
	}

	FormatItems struct {
		Ref string `json:"$ref,omitempty" yaml:"$ref,omitempty"`
	}

	FormatComponents struct {
		Responses map[string]FormatResponse `json:"responses,omitempty" yaml:"responses,omitempty"`
		Schemas   map[string]FormatSchema   `json:"schemas,omitempty" yaml:"schemas,omitempty"`
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

func (f *FormatRoute) AddServer(server FormatServer) {
	f.Servers = append(f.Servers, server)
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
