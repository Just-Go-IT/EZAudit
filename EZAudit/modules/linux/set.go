package linux

import (
	"Just-Go-IT/EZAudit/artifact"
	"Just-Go-IT/EZAudit/global"
	"Just-Go-IT/EZAudit/interact"
	"Just-Go-IT/EZAudit/registry"
)

func init() {
	registry.Register("set", &set{}, false, registry.Linux)
}

type set struct {
	command string
}

func (s set) New(p map[string]interface{}) (global.Module, error) {
	return &s, nil
}

func (s *set) Execute(step *global.Step) (output string, err error) {
	s.command += "set"

	output, err = interact.ShellPipe(s.command, step)
	if err != nil {
		return
	}

	artifact.SaveString(output, *step)
	return
}
