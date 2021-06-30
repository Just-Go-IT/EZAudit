package linux

import (
	"Just-Go-IT/EZAudit/artifact"
	"Just-Go-IT/EZAudit/global"
	"Just-Go-IT/EZAudit/interact"
	"Just-Go-IT/EZAudit/registry"
	"errors"
)

func init() {
	registry.Register("nmcli", &nmcli{}, false, registry.Linux)
}

type nmcli struct {
	terse    bool
	radioAll bool
	command  string
}

func (n nmcli) New(p map[string]interface{}) (global.Module, error) {
	ok := false

	// Checks for optional parameters and if the keys are supported
	for key := range p {
		switch key {
		case "terse":
			n.terse, ok = p["terse"].(bool)
			if !ok {
				return nil, errors.New("the key \"terse\" is set for the module. The value must be a \"bool\"")
			}
		case "radioAll":
			n.radioAll, ok = p["radioAll"].(bool)
			if !ok {
				return nil, errors.New("the key \"radioAll\" is set for the module. The value must be a \"bool\"")
			}
		default:
			if key != "terse" && key != "radioAll" {
				return nil, errors.New("there is no key called: \"" + key + "\" in the module nmcli")
			}
		}
	}

	if (n.terse && n.radioAll) || !(n.terse && n.radioAll) {
		return nil, errors.New("only one Option at the time ist valid")
	}

	return &n, nil
}

func (n *nmcli) Execute(s *global.Step) (output string, err error) {
	n.command += "nmcli "

	if n.radioAll {
		n.command += "radio all"
	}
	if n.terse {
		n.command += "-t connection show "
	}

	output, err = interact.ShellPipe(n.command, s)
	if err != nil {
		return
	}

	artifact.SaveString(output, *s)
	return
}
