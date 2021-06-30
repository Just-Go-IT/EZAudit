// Package linux includes all modules which are usable for linux
package linux

import (
	"Just-Go-IT/EZAudit/artifact"
	"Just-Go-IT/EZAudit/global"
	"Just-Go-IT/EZAudit/interact"
	"Just-Go-IT/EZAudit/registry"
)

func init() {
	registry.Register("alias", &alias{}, false, registry.Linux)
}

type alias struct {
	command string
}

func (a alias) New(p map[string]interface{}) (global.Module, error) {
	return &a, nil
}

func (a *alias) Execute(s *global.Step) (output string, err error) {
	a.command += "alias"

	output, err = interact.ShellPipe(a.command, s)
	if err != nil {
		return
	}

	artifact.SaveString(output, *s)
	return
}
