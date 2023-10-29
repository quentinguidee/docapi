package generator

import (
	"bufio"
	"docapi/types"
	"encoding/json"
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

type apiTemp struct {
	// Groups are the http groups of the API.
	Groups []types.Group
	// Routes are the http routes of the API.
	Routes []types.Route
	// Handlers are the http handlers of the API.
	Handlers map[string]types.Handler

	currentHandler   types.Handler
	currentHandlerID string
}

type (
	api struct {
		Groups []types.Group `json:"groups"`
		Routes []apiRoute    `json:"routes"`
	}

	apiRoute struct {
		types.Route
		types.Handler
	}
)

func newAPI() *apiTemp {
	return &apiTemp{
		Handlers: make(map[string]types.Handler),
	}
}

func (a *apiTemp) Output() (string, error) {
	export := api{
		Groups: a.Groups,
		Routes: make([]apiRoute, 0),
	}
	for _, route := range a.Routes {
		export.Routes = append(export.Routes, apiRoute{
			Route:   route,
			Handler: a.Handlers[route.HandlerID],
		})
	}

	doc, err := json.MarshalIndent(export, "", "  ")
	if err != nil {
		return "", err
	}
	return string(doc), nil
}

func (a *apiTemp) walk(path string) error {
	return filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
		return a.collect(path)
	})
}

func (a *apiTemp) collect(path string) error {
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

func (a *apiTemp) parse(line string) error {
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
	case strings.HasPrefix(line, "status"):
		return a.status(args)
	case strings.HasPrefix(line, "end"):
		return a.end()
	default:
		fmt.Printf("warn: %s: %s\n", ErrInvalidCommand, line)
		return nil
	}
}

func (a *apiTemp) group(args []string) error {
	if len(args) != 1 {
		return ErrInvalidNumberOfArguments
	}
	a.Groups = append(a.Groups, types.Group{
		Path: args[0],
	})
	return nil
}

func (a *apiTemp) route(args []string) error {
	if len(args) != 2 {
		return ErrInvalidNumberOfArguments
	}
	a.Routes = append(a.Routes, types.Route{
		Path:      args[0],
		HandlerID: args[1],
	})
	return nil
}

func (a *apiTemp) begin(args []string) error {
	if len(args) != 1 {
		return ErrInvalidNumberOfArguments
	}
	a.currentHandler = types.Handler{}
	a.currentHandlerID = args[0]
	return nil
}

func (a *apiTemp) method(args []string) error {
	if len(args) != 1 {
		return ErrInvalidNumberOfArguments
	}
	a.currentHandler.Method = args[0]
	return nil
}

func (a *apiTemp) summary(args []string) error {
	if len(args) < 1 {
		return ErrInvalidNumberOfArguments
	}
	a.currentHandler.Summary = strings.Join(args, " ")
	return nil
}

func (a *apiTemp) consumes(args []string) error {
	if len(args) != 1 {
		return ErrInvalidNumberOfArguments
	}
	a.currentHandler.Consumes = args[0]
	return nil
}

func (a *apiTemp) produces(args []string) error {
	if len(args) != 1 {
		return ErrInvalidNumberOfArguments
	}
	a.currentHandler.Produces = args[0]
	return nil
}

func (a *apiTemp) status(args []string) error {
	if len(args) != 1 {
		return ErrInvalidNumberOfArguments
	}
	status, err := strconv.Atoi(args[0])
	if err != nil {
		return err
	}
	a.currentHandler.Status = append(a.currentHandler.Status, status)
	return nil
}

func (a *apiTemp) end() error {
	a.Handlers[a.currentHandlerID] = a.currentHandler
	return nil
}
