package linux

import (
	"Just-Go-IT/EZAudit/artifact"
	"Just-Go-IT/EZAudit/global"
	"Just-Go-IT/EZAudit/interact"
	"Just-Go-IT/EZAudit/registry"
	"errors"
)

func init() {
	registry.Register("sysctl", &sysctl{}, false, registry.Linux)
}

type sysctl struct {
	target  string
	command string
}

func (s sysctl) New(p map[string]interface{}) (global.Module, error) {
	ok := false

	s.target, ok = p["target"].(string)
	if !ok {
		return nil, errors.New("the key \"target\" must be set and the value must be a \"string\"")
	}

	// Checks for optional parameters and if the keys are supported
	for key := range p {
		if key != "target" {
			return nil, errors.New("there is no key called: \"" + key + "\" in the module sysctl")
		}
	}

	return &s, nil
}

func (s *sysctl) Execute(step *global.Step) (output string, err error) {
	s.command += "sysctl " + s.target

	output, err = interact.ShellPipe(s.command, step)

	if err != nil {
		return
	}

	artifact.SaveString(output, *step)
	return
}
