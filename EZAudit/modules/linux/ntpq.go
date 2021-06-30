package linux

import (
	"Just-Go-IT/EZAudit/artifact"
	"Just-Go-IT/EZAudit/global"
	"Just-Go-IT/EZAudit/interact"
	"Just-Go-IT/EZAudit/registry"
)

func init() {
	registry.Register("ntpq", &ntpq{}, false, registry.Linux)
}

type ntpq struct {
	command string
}

func (n ntpq) New(p map[string]interface{}) (global.Module, error) {
	return &n, nil
}

func (n *ntpq) Execute(s *global.Step) (output string, err error) {
	n.command += "ntpq -p -n"

	output, err = interact.ShellPipe(n.command, s)
	if err != nil {
		return
	}

	artifact.SaveString(output, *s)
	return
}
