package linux

import (
	"Just-Go-IT/EZAudit/artifact"
	"Just-Go-IT/EZAudit/global"
	"Just-Go-IT/EZAudit/interact"
	"Just-Go-IT/EZAudit/registry"
)

func init() {
	registry.Register("subscriptionManager", &subscriptionManager{}, false, registry.Linux)
}

type subscriptionManager struct {
	command string
}

func (s subscriptionManager) New(p map[string]interface{}) (global.Module, error) {
	return &s, nil
}

func (s *subscriptionManager) Execute(step *global.Step) (output string, err error) {
	s.command += "subscription-manager identity"

	output, err = interact.ShellPipe(s.command, step)

	if err != nil {
		return
	}

	artifact.SaveString(output, *step)
	return
}
