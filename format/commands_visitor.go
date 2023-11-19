package format

import (
	"fmt"
	"strings"

	"github.com/quentinguidee/docapi/types"
)

type CommandsVisitor struct {
	api *api
}

func NewCommandsVisitor(api *api) *CommandsVisitor {
	return &CommandsVisitor{
		api: api,
	}
}

func (v *CommandsVisitor) Visit(cmd types.Command) error {
	switch cmd.Type {
	case types.CmdTitle:
		v.visitTitle(cmd)
	case types.CmdDescription:
		v.visitDescription(cmd)
	case types.CmdVersion:
		v.visitVersion(cmd)
	case types.CmdFilename:
		v.visitFilename(cmd)
	case types.CmdUrl:
		v.visitUrl(cmd)
	case types.CmdUrlVar:
		v.visitUrlVar(cmd)
	case types.CmdCode:
		v.visitCode(cmd)
	case types.CmdRoute:
		v.visitRoute(cmd)
	case types.CmdBegin:
		v.visitBegin(cmd)
	case types.CmdMethod:
		v.visitMethod(cmd)
	case types.CmdSummary:
		v.visitSummary(cmd)
	case types.CmdDesc:
		v.visitDesc(cmd)
	case types.CmdTags:
		v.visitTags(cmd)
	case types.CmdBody:
		v.visitBody(cmd)
	case types.CmdQuery:
		v.visitQuery(cmd)
	case types.CmdResponse:
		v.visitResponse(cmd)
	case types.CmdEnd:
		v.visitEnd(cmd)
	default:
		return fmt.Errorf("invalid command: %s", cmd.Type)
	}
	return nil
}

func (v *CommandsVisitor) visitTitle(cmd types.Command) {
	v.api.Info.Title = strings.Join(cmd.Args, " ")
}

func (v *CommandsVisitor) visitDescription(cmd types.Command) {
	v.api.Info.Description = strings.Join(cmd.Args, " ")
}

func (v *CommandsVisitor) visitVersion(cmd types.Command) {
	v.api.Info.Version = cmd.Args[0]
}

func (v *CommandsVisitor) visitFilename(cmd types.Command) {
	v.api.filename = cmd.Args[0]
}

func (v *CommandsVisitor) visitUrl(cmd types.Command) {
	v.api.AddServer(types.FormatServer{
		Url: cmd.Args[0],
	})
}

func (v *CommandsVisitor) visitUrlVar(cmd types.Command) {
	var (
		name         = cmd.Args[0]
		defaultValue = cmd.Args[1]
		description  = strings.Join(cmd.Args[2:], " ")
	)
	variable := types.FormatServerVariable{
		Default:     defaultValue,
		Description: description,
	}
	v.api.Servers[len(v.api.Servers)-1].SetVariable(name, variable)
}

func (v *CommandsVisitor) visitCode(cmd types.Command) {
	code := cmd.Args[0]
	args := cmd.Args[1:]
	ref := ""
	resp := types.FormatResponse{}
	if strings.HasPrefix(args[0], "{") {
		ref = args[0][1 : len(args[0])-1]
		args = args[1:]
		resp.Content = map[string]types.FormatContent{
			"application/json": {
				Schema: types.FormatSchema{
					Ref: types.CreateRef(types.RefSchema, ref),
				},
			},
		}
	}
	resp.Description = strings.Join(args, " ")
	v.api.Components.SetResponse(code, resp)
}

func (v *CommandsVisitor) visitRoute(cmd types.Command) {
	v.api.routes[cmd.Args[1]] = cmd.Args[0]
}

func (v *CommandsVisitor) visitBegin(cmd types.Command) {
	v.api.tempHandler = types.FormatRoute{
		OperationId: cmd.Args[0],
	}
}

func (v *CommandsVisitor) visitMethod(cmd types.Command) {
	v.api.handlerMethods[v.api.tempHandler.OperationId] = strings.ToLower(cmd.Args[0])
}

func (v *CommandsVisitor) visitSummary(cmd types.Command) {
	v.api.tempHandler.Summary = strings.Join(cmd.Args, " ")
}

func (v *CommandsVisitor) visitDesc(cmd types.Command) {
	v.api.tempHandler.Description = strings.Join(cmd.Args, " ")
}

func (v *CommandsVisitor) visitTags(cmd types.Command) {
	v.api.tempHandler.Tags = cmd.Args
}

func (v *CommandsVisitor) visitBody(cmd types.Command) {
	component := cmd.Args[0]
	component = component[1 : len(component)-1]
	description := cmd.Args[1:]

	v.api.tempHandler.RequestBody = types.FormatRequestBody{
		Description: strings.Join(description, " "),
		Required:    true,
		Content: map[string]types.FormatContent{
			"application/json": {
				Schema: v.api.schemaFromAlias(component),
			},
		},
	}
}

func (v *CommandsVisitor) visitQuery(cmd types.Command) {
	component := cmd.Args[1]
	component = component[1 : len(component)-1]
	schema := v.api.schemaFromAlias(component)
	v.api.tempHandler.AddParameter(types.FormatParameter{
		In:          "query",
		Name:        cmd.Args[0],
		Description: strings.Join(cmd.Args[2:], " "),
		Required:    true,
		Schema:      schema,
	})
}

func (v *CommandsVisitor) visitResponse(cmd types.Command) {
	if len(cmd.Args) <= 1 {
		v.api.tempHandler.SetResponse(cmd.Args[0], types.FormatResponse{})
		return
	}

	resp := types.FormatResponse{
		Description: strings.Join(cmd.Args[2:], " "),
		Content:     map[string]types.FormatContent{},
	}

	content := types.FormatContent{}
	component := cmd.Args[1]
	component = component[1 : len(component)-1]
	content.Schema = v.api.schemaFromAlias(component)
	resp.Content["application/json"] = content
	v.api.tempHandler.SetResponse(cmd.Args[0], resp)
}

func (v *CommandsVisitor) visitEnd(cmd types.Command) {
	v.api.handlers[v.api.tempHandler.OperationId] = v.api.tempHandler
}
