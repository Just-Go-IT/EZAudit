package linux

import (
	"Just-Go-IT/EZAudit/artifact"
	"Just-Go-IT/EZAudit/global"
	"Just-Go-IT/EZAudit/interact"
	"Just-Go-IT/EZAudit/registry"
	"errors"
)

func init() {
	registry.Register("yum", &yum{}, false, registry.Linux)
}

type yum struct {
	listInstalled bool
	listSecurity  bool
	command       string
}

func (y yum) New(p map[string]interface{}) (global.Module, error) {
	ok := false

	// Checks for optional parameters and if the keys are supported
	for key := range p {
		switch key {
		case "listInstalled":
			y.listInstalled, ok = p["listInstalled"].(bool)
			if !ok {
				return nil, errors.New("the key \"listInstalled\" is set for the module. The value must be a \"bool\"")
			}
		case "listSecurity":
			y.listSecurity, ok = p["listSecurity"].(bool)
			if !ok {
				return nil, errors.New("the key \"listSecurity\" is set for the module. The value must be a \"bool\"")
			}
		default:
			if key != "listinstalled" && key != "listsecurity" {
				return nil, errors.New("there is no key called: \"" + key + "\" in the module yum")
			}
		}
	}

	if (!y.listSecurity && !y.listInstalled) || (y.listInstalled && y.listSecurity) {
		return nil, errors.New("there is no action option in the module yum")
	}

	return &y, nil
}

func (y *yum) Execute(s *global.Step) (output string, err error) {
	y.command += "yum "
	if y.listSecurity {
		y.command += "list-security"
	}
	if y.listSecurity {
		y.command += "list installed"
	}

	output, err = interact.ShellPipe(y.command, s)
	if err != nil {
		return
	}

	artifact.SaveString(output, *s)
	return
}
