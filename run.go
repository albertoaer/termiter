package termiter

import "fmt"

const MaxActionExecution = 50

type ExecutionContext struct {
	visitedActions map[string]int
}

func NewExecutionContext() *ExecutionContext {
	return &ExecutionContext{make(map[string]int)}
}

type Runnable interface {
	Run(context *ExecutionContext, args []string) int
}

func (file *TermiterFile) Action(name string) (Runnable, error) {
	if act, e := file.Actions[name]; e {
		return act, nil
	}
	return nil, fmt.Errorf("Action %s not found", name)
}

func (file *TermiterFile) Command(name string) (Runnable, error) {
	if act, e := file.Commands[name]; e {
		return act, nil
	}
	return nil, fmt.Errorf("Command %s not found", name)
}

func (file *TermiterFile) ActionOrCommand(name string) (Runnable, error) {
	if x, _ := file.Action(name); x != nil {
		return x, nil
	}
	if x, _ := file.Command(name); x != nil {
		return x, nil
	}
	return nil, fmt.Errorf("Neither Command or Action %s found", name)
}

func (file *TermiterFile) GetRunnable(args []string) (Runnable, error) {
	switch len(args) {
	case 0:
		return file.ActionOrCommand("def")
	case 1:
		return file.ActionOrCommand(args[0])
	default:
		return file.Command(args[0])
	}
}

func (action *Action) Run(context *ExecutionContext, _ []string) int {
	//TODO: Implement
	return -1
}

func (command *Command) Run(context *ExecutionContext, args []string) int {
	//TODO: Implement
	return -1
}
