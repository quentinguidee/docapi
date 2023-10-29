package types

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
		Summary string `json:"summary"`
		// Consumes is the content type of the request.
		Consumes string `json:"consumes"`
		// Produces is the content type of the response.
		Produces string `json:"produces"`
		// Statuses are the possible HTTP codes of the response.
		Status []int `json:"status"`
	}
)
