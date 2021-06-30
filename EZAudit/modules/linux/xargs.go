package linux

import (
	"Just-Go-IT/EZAudit/artifact"
	"Just-Go-IT/EZAudit/global"
	"Just-Go-IT/EZAudit/interact"
	"Just-Go-IT/EZAudit/registry"
	"errors"
)

func init() {
	registry.Register("xargs", &xargs{}, false, registry.Linux)
}

type xargs struct {
	null      bool
	arguments string
	command   string
}

func (x xargs) New(p map[string]interface{}) (global.Module, error) {
	ok := false

	// Checks for optional parameters and if the keys are supported
	for key := range p {
		switch key {
		case "null":
			x.null, ok = p["null"].(bool)
			if !ok {
				return nil, errors.New("the key \"null\" is set for the module. The value must be a \"bool\"")
			}
		case "arguments":
			x.arguments, ok = p["arguments"].(string)
			if !ok {
				return nil, errors.New("the key \"arguments\" is set for the module. The value must be a \"string\"")
			}
		default:
			if key != "null" && key != "arguments" {
				return nil, errors.New("there is no key called: \"" + key + "\" in the module xargs")
			}
		}
	}

	if !x.null && x.arguments == "" {
		return nil, errors.New("there is no action option in the module xargs")
	}

	return &x, nil
}

func (x *xargs) Execute(s *global.Step) (output string, err error) {

	x.command += "xargs"

	if x.null {
		x.command += " -0"
	}
	if x.arguments != "" {
		x.command += " " + x.arguments
	}

	output, err = interact.ShellPipe(x.command, s)
	if err != nil {
		return
	}

	artifact.SaveString(output, *s)
	return
}
