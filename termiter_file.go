package termiter

import (
	"fmt"
	"io"

	"github.com/BurntSushi/toml"
)

//A profile is a bridge between the user input and a program execution
type Profile struct {
	Extends string            `toml:"extend"` //The profile from which its inherits its behave
	Accepts []string          `toml:"accept"` //Flags accepted by the profile
	Set     map[string]string `toml:"set"`    //Setted flags by the profile
	Nargs   int               `toml:"nargs"`  //Number of arguments it accepts
}

//An action is an execution that always do the same
type Action struct {
	Expect   []string `toml:"expect"`   //Actions expected to be donde before
	Exec     []string `toml:"exec"`     //Command to be executed
	Once     bool     `toml:"once"`     //If its true and the action has been executed won't run again
	Optional bool     `toml:"optional"` //If its fails the action expecting it will run without problems
}

//A command is a action with custom input provided by the user and mapped by the profile
type Command struct {
	Expect []string `toml:"expect"` //Actions expected to be done before
	Use    string   `toml:"use"`    //Profile to be applied
	Target string   `toml:"target"` //Target program to be run
}

type TermiterFile struct {
	Commands map[string]*Command `toml:"cmd"`
	Actions  map[string]*Action  `toml:"act"`
	Profiles map[string]*Profile `toml:"prf"`
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
	for k := range file.Commands {
		if _, e := file.Actions[k]; e {
			return fmt.Errorf("Repeated key for command and action: %s", k)
		}
	}
	for name, profile := range file.Profiles {
		if _, e := file.Profiles[profile.Extends]; !e && len(profile.Extends) > 0 {
			return fmt.Errorf("Profile %s extending invalid profile: %s", name, profile.Extends)
		}
	}
	for name, action := range file.Actions {
		for _, expect := range action.Expect {
			if _, e := file.Actions[expect]; !e {
				return fmt.Errorf("Action %s expecting invalid action: %s", name, expect)
			}
		}
	}
	for name, command := range file.Commands {
		for _, expect := range command.Expect {
			if _, e := file.Actions[expect]; !e {
				return fmt.Errorf("Command %s expecting invalid action: %s", name, expect)
			}
		}
		if _, e := file.Profiles[command.Use]; !e && len(command.Use) > 0 {
			return fmt.Errorf("Command %s using invalid profile: %s", name, command.Use)
		}
	}
	return nil
}
