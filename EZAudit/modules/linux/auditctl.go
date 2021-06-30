package linux

import (
	"Just-Go-IT/EZAudit/artifact"
	"Just-Go-IT/EZAudit/global"
	"Just-Go-IT/EZAudit/interact"
	"Just-Go-IT/EZAudit/registry"
)

func init() {
	registry.Register("auditctl", &auditctl{}, false, registry.Linux)
}

type auditctl struct {
	command string
}

func (a auditctl) New(p map[string]interface{}) (global.Module, error) {
	return &a, nil
}

func (a *auditctl) Execute(currentStep *global.Step) (output string, err error) {
	a.command += "auditctl-l "

	output, err = interact.ShellPipe(a.command, currentStep)
	if err != nil {
		return
	}

	artifact.SaveString(output, *currentStep)

	return
}
