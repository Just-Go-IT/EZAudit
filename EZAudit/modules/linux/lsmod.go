package linux

import (
	"Just-Go-IT/EZAudit/artifact"
	"Just-Go-IT/EZAudit/global"
	"Just-Go-IT/EZAudit/interact"
	"Just-Go-IT/EZAudit/registry"
)

func init() {
	registry.Register("lsmod", &lsmod{}, false, registry.Linux)
}

type lsmod struct {
	command string
}

func (l lsmod) New(p map[string]interface{}) (global.Module, error) {
	return &l, nil
}

func (l *lsmod) Execute(s *global.Step) (output string, err error) {
	l.command += "lsmod"

	// execute command
	output, err = interact.ShellPipe(l.command, s)
	if err != nil {
		return
	}

	artifact.SaveString(output, *s)

	return
}
