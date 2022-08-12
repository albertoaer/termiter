package termiter

import (
	"fmt"
	"io"

	"github.com/BurntSushi/toml"
)

type Profile struct {
	Extends string            `toml:"extend"`
	Accepts []string          `toml:"accept"`
	Set     map[string]string `toml:"set"`
	Nargs   int               `toml:"nargs"`
}

type Action struct {
	Expect   []string `toml:"expect"`
	Exec     []string `toml:"exec"`
	Once     bool     `toml:"once"`
	Optional bool     `toml:"optional"`
}

type Command struct {
	Expect []string `toml:"expect"`
	Use    string   `toml:"use"`
	Target string   `toml:"target"`
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
