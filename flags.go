package termiter

import (
	"gopkg.in/alecthomas/kingpin.v2"
)

func getFlagValues(flags map[string]*Flag, args []string) (map[string]string, error) {
	output := make(map[string]string)
	arr := make(map[string]*string)
	app := kingpin.New("Termiter", "Termiter build tool")
	for k, v := range flags {
		flag := app.Flag(v.Name, v.Info)
		flag.Default(v.Default)
		if v.Required {
			flag.Required()
		}
		arr[k] = flag.String()
	}
	if _, e := app.Parse(args); e != nil {
		return nil, e
	}
	for k, v := range arr {
		output[k] = *v
	}
	return output, nil
}
