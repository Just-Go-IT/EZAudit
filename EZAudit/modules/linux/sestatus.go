package linux

import (
	"Just-Go-IT/EZAudit/artifact"
	"Just-Go-IT/EZAudit/global"
	"Just-Go-IT/EZAudit/interact"
	"Just-Go-IT/EZAudit/registry"
	"errors"
)

func init() {
	registry.Register("sestatus", &sestatus{}, false, registry.Linux)
}

type sestatus struct {
	optionV bool
	command string
}

func (s sestatus) New(p map[string]interface{}) (global.Module, error) {
	ok := false

	// Checks for optional parameters and if the keys are supported
	for key := range p {
		switch key {
		case "optionV":
			s.optionV, ok = p["optionV"].(bool)
			if !ok {
				return nil, errors.New("the key \"optionV\" is set for the module. The value must be a \"bool\"")
			}
		default:
			if key != "optionV" {
				return nil, errors.New("there is no key called: \"" + key + "\" in the module sestatus")
			}
		}
	}

	if !s.optionV {
		return nil, errors.New("there is no action option in the module sestatus")
	}

	return &s, nil
}

func (s *sestatus) Execute(step *global.Step) (output string, err error) {
	s.command += "sestatus "

	if s.optionV {
		s.command += "-v "
	}

	output, err = interact.ShellPipe(s.command, step)

	if err != nil {
		return
	}

	artifact.SaveString(output, *step)
	return
}
