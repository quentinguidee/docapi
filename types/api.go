package types

type (
	Method string

	Api struct {
		Groups []Group                        `json:"groups"`
		Routes map[string]map[Method]ApiRoute `json:"routes"`
	}

	ApiRoute struct {
		Route
		Handler
	}
)

type (
	// Group describes an http group.
	Group struct {
		// Path is the path of the group.
		Path string `json:"path"`
	}

	// Route describes an http route.
	Route struct {
		// Path is the path of the route.
		Path string `json:"path"`
		// ID is the unique identifier of the handler.
		HandlerID string `json:"-"`
	}

	// Handler describes an http handler.
	Handler struct {
		// Method is the HTTP method of the handler.
		Method string `json:"method"`
		// Summary is description of the handler.
		Summary string `json:"summary,omitempty"`
		// Consumes is the content type of the request.
		Consumes string `json:"consumes,omitempty"`
		// Produces is the content type of the response.
		Produces string `json:"produces,omitempty"`
		// Statuses are the possible HTTP codes of the response.
		Responses []Response `json:"responses"`
	}

	Response struct {
		// Code is the HTTP code of the response.
		Code int `json:"code"`
		// Type is the type of the response.
		Type string `json:"type,omitempty"`
		// Value is the value of the response.
		Ref string `json:"ref,omitempty"`
	}
)
