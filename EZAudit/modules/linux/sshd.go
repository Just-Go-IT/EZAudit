package linux

import (
	"Just-Go-IT/EZAudit/artifact"
	"Just-Go-IT/EZAudit/global"
	"Just-Go-IT/EZAudit/interact"
	"Just-Go-IT/EZAudit/registry"
)

func init() {
	registry.Register("sshd", &sshd{}, false, registry.Linux)
}

type sshd struct {
	command string
}

func (s sshd) New(p map[string]interface{}) (global.Module, error) {
	return &s, nil
}

func (s *sshd) Execute(step *global.Step) (output string, err error) {
	s.command += "sshd -T "

	output, err = interact.ShellPipe(s.command, step)

	if err != nil {
		return
	}

	artifact.SaveString(output, *step)
	return
}
