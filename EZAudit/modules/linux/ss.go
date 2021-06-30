package linux

import (
	"Just-Go-IT/EZAudit/artifact"
	"Just-Go-IT/EZAudit/global"
	"Just-Go-IT/EZAudit/interact"
	"Just-Go-IT/EZAudit/registry"
	"errors"
)

func init() {
	registry.Register("ss", &ss{}, false, registry.Linux)
}

type ss struct {
	argument string
	command  string
}

func (s ss) New(p map[string]interface{}) (global.Module, error) {
	ok := false

	// Checks for optional parameters and if the keys are supported
	for key := range p {
		switch key {
		case "argument":
			s.argument, ok = p["argument"].(string)
			if !ok {
				return nil, errors.New("the key \"argument\" is set for the module. The value must be a \"string\"")
			}
		default:
			if key != "argument" {
				return nil, errors.New("there is no key called: \"" + key + "\" in the module mount")
			}
		}
	}

	return &s, nil
}

func (s *ss) Execute(step *global.Step) (output string, err error) {
	s.command += "ss "

	if s.argument != "" {
		s.command += s.argument
	}

	output, err = interact.ShellPipe(s.command, step)

	if err != nil {
		return
	}

	artifact.SaveString(output, *step)
	return
}
