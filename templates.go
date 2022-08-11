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
	Exec   []string `toml:"exec"`
	Use    string   `toml:"use"`
}

type TermiterFile struct {
	Commands map[string]Command `toml:"cmd"`
	Actions  map[string]Action  `toml:"act"`
	Profiles map[string]Profile `toml:"prf"`
}

func ReadTermiterFile(source io.Reader) (file TermiterFile, err error) {
	decoder := toml.NewDecoder(source)
	_, err = decoder.Decode(&file)
	return
}

func Verify(file TermiterFile) error {
	for k := range file.Commands {
		if _, e := file.Actions[k]; e {
			return fmt.Errorf("Repeated key for command and action: %s", k)
		}
	}
	return nil
}
