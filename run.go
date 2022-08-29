package termiter

import (
	"os"
)

const MaxActionExecution = 50

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

func makeEnv(name string, context *ExecutionContext) ([]string, error) {
	if env, e := context.file.Environments[name]; e {
		result := make([]string, 0)
		if env.Fork {
			result = append(result, os.Environ()...)
		}
		for k, v := range env.Include {
			var err error
			if k, err = replaceVariable(context, k); err != nil {
				return nil, err
			}
			if v, err = replaceVariable(context, v); err != nil {
				return nil, err
			}
			result = append(result, k+"="+v)
		}
		return result, nil
	}
	return os.Environ(), nil
}

func (action *Action) Run(context *ExecutionContext) int {
	RunExpected(context, action.Expect)
	if len(action.Exec) > 0 {
		env, err1 := makeEnv(action.Env, context)
		cmd, err2 := parseCommand(action.Exec, context)
		PanicIfError(err1, err2)
		return RunCommand(cmd[0], cmd[1:], env)
	}
	return 0
}
