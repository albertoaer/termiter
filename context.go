package termiter

import (
	"errors"
	"fmt"

	"gopkg.in/alecthomas/kingpin.v2"
)

type ExecutionContext struct {
	file           *TermiterFile
	visitedActions map[string]int
	flags          map[string]*string
	action         string
}

func NewExecutionContext(file *TermiterFile, args []string) (*ExecutionContext, error) {
	app := kingpin.New("Termiter", "Termiter build tool")
	flagsPtr := includeFlags(app, file.Flags)
	includeActions(app, file.Actions)
	var err error
	var action string
	if action, err = app.Parse(args); err != nil {
		return nil, err
	}
	return &ExecutionContext{file, make(map[string]int), flagsPtr, action}, nil
}

func (context *ExecutionContext) GetVariable(name string) (value string, err error) {
	if v, e := context.file.Variables[name]; e {
		return v.Value, nil
	}
	if v, e := context.flags[name]; e {
		return *v, nil
	}
	return "", fmt.Errorf("Variable not found: '%s'", name)
}

func (context *ExecutionContext) GetStartAction() (value *Action, err error) {
	ac := context.file.Actions[context.action]
	if ac != nil {
		return ac, nil
	}
	return nil, errors.New("Invalid valid start action")
}

func includeActions(app *kingpin.Application, actions map[string]*Action) {
	for k, v := range actions {
		cmd := app.Command(k, v.Info)
		if k == "def" {
			cmd.Default()
		}
	}
}

func includeFlags(app *kingpin.Application, flags map[string]*Flag) (output map[string]*string) {
	output = make(map[string]*string)
	for k, v := range flags {
		flag := app.Flag(v.Name, v.Info)
		flag.Default(v.Default)
		if v.Required {
			flag.Required()
		}
		output[k] = flag.String()
	}
	return output
}
