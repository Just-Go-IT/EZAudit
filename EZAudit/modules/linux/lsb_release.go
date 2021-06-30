package linux

import (
	"Just-Go-IT/EZAudit/artifact"
	"Just-Go-IT/EZAudit/global"
	"Just-Go-IT/EZAudit/interact"
	"Just-Go-IT/EZAudit/registry"
)

func init() {
	registry.Register("lsbRelease", &lsbRelease{}, false, registry.Linux)
}

type lsbRelease struct {
	command string
}

func (l lsbRelease) New(p map[string]interface{}) (global.Module, error) {
	return &l, nil
}

func (l *lsbRelease) Execute(s *global.Step) (output string, err error) {
	l.command += "lsb_release -i"

	output, err = interact.ShellPipe(l.command, s)
	if err != nil {
		return
	}

	artifact.SaveString(output, *s)

	return
}
