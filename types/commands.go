package types

type CommandType string

var (
	CmdTitle       CommandType = "title"
	CmdDescription CommandType = "description"
	CmdVersion     CommandType = "version"
	CmdGroup       CommandType = "group"
	CmdRoute       CommandType = "route"
	CmdBegin       CommandType = "begin"
	CmdMethod      CommandType = "method"
	CmdSummary     CommandType = "summary"
	CmdBody        CommandType = "body"
	CmdResponse    CommandType = "response"
	CmdEnd         CommandType = "end"
)

type Command struct {
	Type CommandType
	Args []string
}