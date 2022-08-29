package termiter

import (
	"errors"
	"strings"
)

func replaceVariable(context *ExecutionContext, line string) (string, error) {
	result := ""
	remain := line
	for len(remain) > 0 {
		if i := strings.IndexAny(remain, "{}"); i >= 0 {
			if remain[i] == '}' {
				return "", errors.New("Found close bracket before open bracket")
			}
			j := strings.IndexAny(remain[i+1:], "{}")
			if j < 0 {
				return "", errors.New("Expecting close bracket")
			}
			if remain[i+1+j] == '{' {
				return "", errors.New("Fount open bracket when expecting closing bracket")
			}
			name := remain[i+1 : i+1+j]
			v, e := context.GetVariable(name)
			if e != nil {
				return "", e
			}
			result += remain[:i] + v + remain[i+j+2:]
			remain = remain[i+j+2:]
		} else {
			result += remain
			remain = ""
		}
	}
	return result, nil
}

func parseCommand(input []string, context *ExecutionContext) (output []string, err error) {
	for _, v := range input {
		if len(v) != 0 && v[0] == '$' {
			output = append(output, strings.Split(v[1:], " ")...)
		} else {
			output = append(output, v)
		}
	}
	for i, e := range output {
		output[i], err = replaceVariable(context, e)
		if err != nil {
			return
		}
	}
	return
}
