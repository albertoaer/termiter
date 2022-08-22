package termiter

import "fmt"

const MaxActionExecution = 50

type ExecutionContext struct {
	file           *TermiterFile
	visitedActions map[string]int
}

func NewExecutionContext(file *TermiterFile, args []string) *ExecutionContext {
	return &ExecutionContext{file, make(map[string]int)}
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

func (action *Action) Run(context *ExecutionContext) int {
	RunExpected(context, action.Expect)
	if len(action.Exec) > 0 {
		return RunCommand(action.Exec[0], action.Exec[1:])
	}
	return 0
}
