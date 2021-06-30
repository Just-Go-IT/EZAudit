package linux

import (
	"Just-Go-IT/EZAudit/artifact"
	"Just-Go-IT/EZAudit/global"
	"Just-Go-IT/EZAudit/interact"
	"Just-Go-IT/EZAudit/registry"
	"errors"
)

func init() {
	registry.Register("stat", &stat{}, false, registry.Linux)
}

type stat struct {
	path    string
	command string
}

func (s stat) New(p map[string]interface{}) (global.Module, error) {
	ok := false

	s.path, ok = p["path"].(string)
	if !ok {
		return nil, errors.New("the key \"path\" must be set and the value must be a \"string\"")
	}

	// Checks for optional parameters and if the keys are supported
	for key := range p {
		if key != "path" {
			return nil, errors.New("there is no key called: \"" + key + "\" in the module sysctl")
		}
	}

	return &s, nil
}

func (s *stat) Execute(step *global.Step) (output string, err error) {
	s.command += "stat " + "\"" + s.path + "\""

	output, err = interact.ShellPipe(s.command, step)

	if err != nil {
		return
	}

	artifact.SaveString(output, *step)
	return
}
