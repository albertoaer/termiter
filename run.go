package termiter

import (
	"fmt"
	"os"
)

const MaxActionExecution = 50

type ExecutionContext struct {
	file           *TermiterFile
	visitedActions map[string]int
	variables      map[string]string
}

func NewExecutionContext(file *TermiterFile, args []string) *ExecutionContext {
	vars := make(map[string]string)
	for k, v := range file.Variables {
		vars[k] = v.Value
	}
	return &ExecutionContext{file, make(map[string]int), vars}
}

type Runnable interface {
	Run(context *ExecutionContext) int
}

func (file *TermiterFile) GetRunnable(args []string) (Runnable, error) {
	name := "def"
	if len(args) > 0 {
		name = args[0]
	}
	if act, e := file.Actions[name]; e {
		return act, nil
	}
	return nil, fmt.Errorf("Action %s not found", name)
}

func RunExpected(context *ExecutionContext, expectedList []string) int {
	for _, item := range expectedList {
		ac := context.file.Actions[item]
		if num, e := context.visitedActions[item]; !e || (num < MaxActionExecution && !ac.Once) {
			context.visitedActions[item]++
			//TODO: Log status errors and omit optionals
			if status := ac.Run(context); status != 0 {
				return status
			}
		}
	}
	return 0
}

func makeEnv(name string, context *ExecutionContext) []string {
	if env, e := context.file.Environments[name]; e {
		result := make([]string, 0)
		if env.Fork {
			result = append(result, os.Environ()...)
		}
		for k, v := range env.Include {
			result = append(result, k+"="+v)
		}
		return result
	}
	return os.Environ()
}

func (action *Action) Run(context *ExecutionContext) int {
	RunExpected(context, action.Expect)
	if len(action.Exec) > 0 {
		env := makeEnv(action.Env, context)
		cmd, err := parseCommand(action.Exec, context)
		PanicIfError(err)
		return RunCommand(cmd[0], cmd[1:], env)
	}
	return 0
}
