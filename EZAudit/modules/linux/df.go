package linux

import (
	"Just-Go-IT/EZAudit/artifact"
	"Just-Go-IT/EZAudit/global"
	"Just-Go-IT/EZAudit/interact"
	"Just-Go-IT/EZAudit/registry"
	"errors"
)

func init() {
	registry.Register("df", &df{}, false, registry.Linux)
}

type df struct {
	local   bool
	path    string
	command string
}

func (d df) New(p map[string]interface{}) (global.Module, error) {
	ok := false
	d.path, ok = p["path"].(string)
	if !ok {
		return nil, errors.New("the key \"path\" must be set and the value must be a \"string\"")
	}

	// Checks for optional parameter Mode and check if the keys are allowed
	for key := range p {
		switch key {
		case "local":
			d.local, ok = p["local"].(bool)
			if !ok {
				return nil, errors.New("the key \"local\" is set for the module. The value must be a \"bool\"")
			}
		default:
			if key != "local" && key != "path" {
				return nil, errors.New("there is no key called: \"" + key + "\" in the module df")
			}
		}
	}

	return &d, nil
}

func (d *df) Execute(s *global.Step) (output string, err error) {
	d.command = "df"

	if d.local {
		d.command += " -l"
	}
	if d.path != "" {
		d.command += " " + d.path
	}

	// execute command
	output, err = interact.ShellPipe(d.command, s)

	if err != nil {
		return
	}

	artifact.SaveString(output, *s)

	return
}
