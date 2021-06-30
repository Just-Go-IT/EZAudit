package linux

import (
	"Just-Go-IT/EZAudit/artifact"
	"Just-Go-IT/EZAudit/global"
	"Just-Go-IT/EZAudit/interact"
	"Just-Go-IT/EZAudit/registry"
)

func init() {
	registry.Register("lastlog", &lastlog{}, false, registry.Linux)
}

type lastlog struct {
	command string
}

func (l lastlog) New(p map[string]interface{}) (global.Module, error) {
	return &l, nil
}

func (l *lastlog) Execute(s *global.Step) (output string, err error) {
	l.command += "lastlog"

	output, err = interact.ShellPipe(l.command, s)
	if err != nil {
		return
	}

	artifact.SaveString(output, *s)

	return
}
