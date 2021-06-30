package linux

import (
	"Just-Go-IT/EZAudit/artifact"
	"Just-Go-IT/EZAudit/global"
	"Just-Go-IT/EZAudit/interact"
	"Just-Go-IT/EZAudit/registry"
	"errors"
)

func init() {
	registry.Register("journalctl", &journalctl{}, false, registry.Linux)
}

type journalctl struct {
	admin   bool
	command string
}

func (j journalctl) New(p map[string]interface{}) (global.Module, error) {
	ok := false
	// Checks for optional parameters and if the keys are supported
	for key := range p {
		switch key {
		case "admin":
			j.admin, ok = p["admin"].(bool)
			if !ok {
				return nil, errors.New("the key \"admin\" is set for the module. The value must be a \"bool\"")
			}
		default:
			if key != "admin" {
				return nil, errors.New("there is no key called: \"" + key + "\" in the module journalctl")
			}
		}
	}
	return &j, nil
}

func (j *journalctl) Execute(s *global.Step) (output string, err error) {
	if j.admin {
		j.command += "sudo"
	}
	j.command += " journalctl"

	// execute command
	output, err = interact.ShellPipe(j.command, s)
	if err != nil {
		return
	}

	artifact.SaveString(output, *s)
	return
}
