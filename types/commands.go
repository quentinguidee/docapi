package types

type CommandType string

var (
	CmdTitle       CommandType = "title"
	CmdDescription CommandType = "description"
	CmdVersion     CommandType = "version"
	CmdFilename    CommandType = "filename"
	CmdUrl         CommandType = "url"
	CmdUrlVar      CommandType = "urlvar"
	CmdCode        CommandType = "code"
	CmdRoute       CommandType = "route"
	CmdBegin       CommandType = "begin"
	CmdMethod      CommandType = "method"
	CmdSummary     CommandType = "summary"
	CmdDesc        CommandType = "desc"
	CmdTags        CommandType = "tags"
	CmdBody        CommandType = "body"
	CmdQuery       CommandType = "query"
	CmdResponse    CommandType = "response"
	CmdEnd         CommandType = "end"
)

type CommandsVisitor interface {
	Visit(cmd Command) error

	visitTitle(cmd Command)
	visitDescription(cmd Command)
	visitVersion(cmd Command)
	visitFilename(cmd Command)
	visitUrl(cmd Command)
	visitUrlVar(cmd Command)
	visitCode(cmd Command)
	visitRoute(cmd Command)
	visitBegin(cmd Command)
	visitMethod(cmd Command)
	visitSummary(cmd Command)
	visitTags(cmd Command)
	visitBody(cmd Command)
	visitQuery(cmd Command)
	visitResponse(cmd Command)
	visitEnd(cmd Command)
}

type Command struct {
	Type CommandType
	Args []string

	// ServerAlias allows executing this command only for a specific server.
	ServerAlias string
}
