package collector

import (
	"bufio"
	"docapi/types"
	"errors"
	"os"
	"path/filepath"
	"strings"
)

var (
	ErrInvalidNumberOfArguments = errors.New("invalid number of arguments")
	ErrInvalidCommand           = errors.New("invalid command")
)

type CommandsCollector struct {
	Commands []types.Command
}

func NewCommandsCollector() *CommandsCollector {
	return &CommandsCollector{}
}

func (a *CommandsCollector) Run(path string) ([]types.Command, error) {
	err := filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
		return a.collect(path)
	})
	if err != nil {
		return nil, err
	}
	return a.Commands, nil
}

func (a *CommandsCollector) collect(path string) error {
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

func (a *CommandsCollector) parse(line string) error {
	line = strings.TrimSpace(line)

	if !strings.HasPrefix(line, "// docapi:") {
		return nil
	}

	line = strings.TrimSpace(line[10:])
	args := strings.Split(line, " ")
	a.Commands = append(a.Commands, types.Command{
		Type: types.CommandType(args[0]),
		Args: args[1:],
	})
	return nil
}
