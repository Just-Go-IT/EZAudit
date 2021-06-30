package linux

import (
	"Just-Go-IT/EZAudit/artifact"
	"Just-Go-IT/EZAudit/global"
	"Just-Go-IT/EZAudit/interact"
	"Just-Go-IT/EZAudit/registry"
	"errors"
)

func init() {
	registry.Register("find", &find{}, false, registry.Linux)
}

type find struct {
	target  string
	command string
}

func (f find) New(p map[string]interface{}) (global.Module, error) {
	ok := false
	f.target, ok = p["target"].(string)
	if !ok {
		return nil, errors.New("the key \"target\" must be set and the value must be a \"string\"")
	}

	// Checks for optional parameter Mode and check if the keys are allowed
	for key := range p {
		if key != "target" {
			return nil, errors.New("there is no key called: \"" + key + "\" in the module find")
		}
	}

	return &f, nil
}

func (f *find) Execute(s *global.Step) (output string, err error) {
	f.command = "find "

	if f.target != "" {
		f.command += f.target
	}

	output, err = interact.ShellPipe(f.command, s)

	if err != nil {
		return
	}

	artifact.SaveString(output, *s)

	return
}
