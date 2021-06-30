package linux

import (
	"Just-Go-IT/EZAudit/artifact"
	"Just-Go-IT/EZAudit/global"
	"Just-Go-IT/EZAudit/interact"
	"Just-Go-IT/EZAudit/registry"
)

func init() {
	registry.Register("authselectCurrent", &authselectCurrent{}, false, registry.Linux)
}

type authselectCurrent struct {
	command string
}

func (a authselectCurrent) New(p map[string]interface{}) (global.Module, error) {
	return &a, nil
}

func (a *authselectCurrent) Execute(s *global.Step) (output string, err error) {
	a.command += "authselect current "

	// execute command
	output, err = interact.ShellPipe(a.command, s)

	if err != nil {
		return
	}
	artifact.SaveString(output, *s)

	return
}
