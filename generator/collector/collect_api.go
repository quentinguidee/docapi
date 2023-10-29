package collector

import (
	"bufio"
	"docapi/types"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

var (
	ErrInvalidNumberOfArguments = errors.New("invalid number of arguments")
	ErrInvalidCommand           = errors.New("invalid command")
)

type ApiCollector struct {
	// Groups are the http groups of the API.
	Groups []types.Group
	// Routes are the http routes of the API.
	Routes []types.Route
	// Handlers are the http handlers of the API.
	Handlers map[string]types.Handler

	currentHandler   types.Handler
	currentHandlerID string
}

func NewAPICollector() *ApiCollector {
	return &ApiCollector{
		Handlers: make(map[string]types.Handler),
	}
}

func (a *ApiCollector) Output() (interface{}, error) {
	api := types.Api{
		Groups: a.Groups,
		Routes: map[string]map[types.Method]types.ApiRoute{},
	}
	for _, route := range a.Routes {
		r := types.ApiRoute{
			Route:   route,
			Handler: a.Handlers[route.HandlerID],
		}
		if api.Routes[route.Path] == nil {
			api.Routes[route.Path] = map[types.Method]types.ApiRoute{}
		}
		api.Routes[route.Path][types.Method(r.Handler.Method)] = r
	}
	return api, nil
}

func (a *ApiCollector) Run(path string) error {
	return filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
		return a.collect(path)
	})
}

func (a *ApiCollector) collect(path string) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		err := a.parse(line)
		if err != nil {
			return err
		}
	}
	return nil
}

func (a *ApiCollector) parse(line string) error {
	line = strings.TrimSpace(line)

	if !strings.HasPrefix(line, "// docapi:") {
		return nil
	}

	line = strings.TrimSpace(line[10:])
	args := strings.Split(line, " ")[1:]

	switch {
	// Groups
	case strings.HasPrefix(line, "group"):
		return a.group(args)
	// Routes
	case strings.HasPrefix(line, "route"):
		return a.route(args)
	// Handlers
	case strings.HasPrefix(line, "begin"):
		return a.begin(args)
	case strings.HasPrefix(line, "method"):
		return a.method(args)
	case strings.HasPrefix(line, "summary"):
		return a.summary(args)
	case strings.HasPrefix(line, "consumes"):
		return a.consumes(args)
	case strings.HasPrefix(line, "produces"):
		return a.produces(args)
	case strings.HasPrefix(line, "response"):
		return a.response(args)
	case strings.HasPrefix(line, "end"):
		return a.end()
	default:
		fmt.Printf("warn: %s: %s\n", ErrInvalidCommand, line)
		return nil
	}
}

func (a *ApiCollector) group(args []string) error {
	if len(args) != 1 {
		return ErrInvalidNumberOfArguments
	}
	a.Groups = append(a.Groups, types.Group{
		Path: args[0],
	})
	return nil
}

func (a *ApiCollector) route(args []string) error {
	if len(args) != 2 {
		return ErrInvalidNumberOfArguments
	}
	a.Routes = append(a.Routes, types.Route{
		Path:      args[0],
		HandlerID: args[1],
	})
	return nil
}

func (a *ApiCollector) begin(args []string) error {
	if len(args) != 1 {
		return ErrInvalidNumberOfArguments
	}
	a.currentHandler = types.Handler{}
	a.currentHandlerID = args[0]
	return nil
}

func (a *ApiCollector) method(args []string) error {
	if len(args) != 1 {
		return ErrInvalidNumberOfArguments
	}
	a.currentHandler.Method = args[0]
	return nil
}

func (a *ApiCollector) summary(args []string) error {
	if len(args) < 1 {
		return ErrInvalidNumberOfArguments
	}
	a.currentHandler.Summary = strings.Join(args, " ")
	return nil
}

func (a *ApiCollector) consumes(args []string) error {
	if len(args) != 1 {
		return ErrInvalidNumberOfArguments
	}
	a.currentHandler.Consumes = args[0]
	return nil
}

func (a *ApiCollector) produces(args []string) error {
	if len(args) != 1 {
		return ErrInvalidNumberOfArguments
	}
	a.currentHandler.Produces = args[0]
	return nil
}

func (a *ApiCollector) response(args []string) error {
	if len(args) < 1 {
		return ErrInvalidNumberOfArguments
	}
	status, err := strconv.Atoi(args[0])
	if err != nil {
		return err
	}
	res := types.Response{
		Code: status,
	}
	if len(args) > 1 {
		t := args[1]
		if strings.HasPrefix(t, "[]") {
			res.Type = "array"
			res.Ref = t[2:]
		} else {
			res.Type = t
		}
	}
	a.currentHandler.Responses = append(a.currentHandler.Responses, res)
	return nil
}

func (a *ApiCollector) end() error {
	a.Handlers[a.currentHandlerID] = a.currentHandler
	return nil
}
