package types

type Value struct {
	// Type is the type of the struct.
	Type string `json:"type"`
	// Fields are the fields of the struct.
	Fields map[string]interface{} `json:"fields,omitempty"`
}
