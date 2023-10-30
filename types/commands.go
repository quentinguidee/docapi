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
	CmdTags        CommandType = "tags"
	CmdBody        CommandType = "body"
	CmdQuery       CommandType = "query"
	CmdResponse    CommandType = "response"
	CmdEnd         CommandType = "end"
)

type Command struct {
	Type CommandType
	Args []string

	// ServerAlias allows executing this command only for a specific server.
	ServerAlias string
}
