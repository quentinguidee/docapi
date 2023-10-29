package types

type (
	Format struct {
		Openapi    string                  `json:"openapi" yaml:"openapi"`
		Info       FormatInfo              `json:"info" yaml:"info"`
		Paths      map[string]FormatRoutes `json:"paths,omitempty" yaml:"paths,omitempty"`
		Components FormatComponents        `json:"components,omitempty" yaml:"components,omitempty"`
	}

	FormatInfo struct {
		Title   string `json:"title" yaml:"title"`
		Desc    string `json:"description" yaml:"description"`
		Version string `json:"version" yaml:"version"`
	}

	FormatRoutes map[string]FormatRoute

	FormatRoute struct {
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
