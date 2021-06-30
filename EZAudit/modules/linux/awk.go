package linux

import (
	"Just-Go-IT/EZAudit/artifact"
	"Just-Go-IT/EZAudit/global"
	"Just-Go-IT/EZAudit/interact"
	"Just-Go-IT/EZAudit/registry"
	"errors"
)

func init() {
	registry.Register("awk", &awk{}, false, registry.Linux)
}

type awk struct {
	fieldSeparator bool
	pattern        string
	command        string
}

func (a awk) New(p map[string]interface{}) (global.Module, error) {
	ok := false
	a.pattern, ok = p["pattern"].(string)
	if !ok {
		return nil, errors.New("the key \"pattern\" must be set and the value must be a \"string\"")
	}

	// Checks for optional parameters and if the keys are supported
	for key, _ := range p {
		switch key {
		case "fieldSeparator":
			a.fieldSeparator, ok = p["fieldSeparator"].(bool)
			if !ok {
				return nil, errors.New("the key \"fieldSeparator\" is set for the module. The value must be a \"bool\"")
			}
		default:
			if key != "pattern" && key != "fieldSeparator" {
				return nil, errors.New("there is no key called: \"" + key + "\" in the module dpkg")
			}
		}
	}
	if !a.fieldSeparator && a.pattern == "" {
		return nil, errors.New("there is no action option in the module awk")
	}

	return &a, nil
}

func (a *awk) Execute(s *global.Step) (output string, err error) {
	a.command = "awk "
	if a.fieldSeparator {
		a.command += "-F "
	}
	a.command += a.pattern

	output, err = interact.ShellPipe(a.command, s)
	if err != nil {
		return
	}

	artifact.SaveString(output, *s)

	return
}
