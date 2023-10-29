package openapi

type (
	Format struct {
		Paths      map[string]FormatRoutes `json:"paths,omitempty" yaml:"paths,omitempty"`
		Components FormatComponents        `json:"components,omitempty" yaml:"components,omitempty"`
	}

	FormatRoutes map[string]FormatRoute

	FormatRoute struct {
		Summary     string                 `json:"summary,omitempty" yaml:"summary,omitempty"`
		Description string                 `json:"description,omitempty" yaml:"description,omitempty"`
		Responses   map[int]FormatResponse `json:"responses,omitempty" yaml:"responses,omitempty"`
	}

	FormatResponse struct {
		Description string                   `json:"description,omitempty" yaml:"description,omitempty"`
		Content     map[string]FormatContent `json:"content,omitempty" yaml:"content,omitempty"`
	}

	FormatContent struct {
		Schema FormatSchema `json:"schema,omitempty" yaml:"schema,omitempty"`
	}

	FormatSchema struct {
		Type       string                  `json:"type" yaml:"type"`
		Items      []FormatSchema          `json:"items,omitempty" yaml:"items,omitempty"`
		Properties map[string]FormatSchema `json:"properties,omitempty" yaml:"properties,omitempty"`
		Ref        string                  `json:"$ref,omitempty" yaml:"$ref,omitempty"`
	}

	FormatComponents struct {
		Schemas map[string]FormatSchema `json:"schemas,omitempty" yaml:"schemas,omitempty"`
	}
)
