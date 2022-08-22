package termiter

import (
	"fmt"
	"io"

	"github.com/BurntSushi/toml"
)

type Flag struct {
	Name string `toml:"name"` //Expected CLI flag name
}

type Environment struct {
	Fork    bool              `toml:"fork"`    //Does clone from the process environment
	Include map[string]string `toml:"include"` //The key=value to be included
}

type Variable struct {
	Value string `toml:"value"` //Value of the variable
}

//An action is an execution that always do the same
type Action struct {
	Expect   []string `toml:"expect"`   //Actions expected to be donde before
	Exec     []string `toml:"exec"`     //Command to be executed
	Once     bool     `toml:"once"`     //If its true and the action has been executed won't run again
	Optional bool     `toml:"optional"` //If its fails the action expecting it will run without problems
	Env      string   `toml:"env"`      //Environment to use in the process invokation
}

type TermiterFile struct {
	Flags        map[string]*Flag        `toml:"flg"`
	Environments map[string]*Environment `toml:"env"`
	Variables    map[string]*Variable    `toml:"var"`
	Actions      map[string]*Action      `toml:"act"`
}

func ReadTermiterFile(source io.Reader) (file *TermiterFile, err error) {
	decoder := toml.NewDecoder(source)
	file = &TermiterFile{}
	var meta toml.MetaData
	meta, err = decoder.Decode(file)
	if err == nil && len(meta.Undecoded()) > 0 {
		err = fmt.Errorf("Unexpected field: %s", meta.Undecoded()[0].String())
	}
	return
}

func (file *TermiterFile) Verify() error {
	for name, action := range file.Actions {
		for _, expect := range action.Expect {
			if _, e := file.Actions[expect]; !e {
				return fmt.Errorf("Action %s expecting invalid action: %s", name, expect)
			}
		}
		if _, e := file.Environments[action.Env]; !e && len(action.Env) > 0 {
			return fmt.Errorf("Action %s using invalid environment: %s", name, action.Env)
		}
	}
	return nil
}
