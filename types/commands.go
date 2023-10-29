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
	CmdTags        CommandType = "tags"
	CmdBody        CommandType = "body"
	CmdQuery       CommandType = "query"
	CmdResponse    CommandType = "response"
	CmdEnd         CommandType = "end"
)

type Command struct {
	Type CommandType
	Args []string
}
