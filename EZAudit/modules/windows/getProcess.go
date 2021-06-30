package windows

import (
	"Just-Go-IT/EZAudit/artifact"
	"Just-Go-IT/EZAudit/global"
	"Just-Go-IT/EZAudit/interact"
	"Just-Go-IT/EZAudit/registry"
	"errors"
)

func init() {
	registry.Register("getProcess", &getProcess{}, false, registry.Windows)
}

type getProcess struct {
	name    string
	command string
}

func (ps getProcess) New(p map[string]interface{}) (global.Module, error) {
	// Parse arguments
	ok := true
	// Check for optional parameter Mode and whitelist the keys
	for key, _ := range p {
		switch key {
		case "name":
			ps.name, ok = p["name"].(string)
			if !ok {
				return nil, errors.New("the key \"name\" is set for the module. The value must be a \"string\"")
			}
		default:
			if key != "name" {
				return nil, errors.New("there is no key called: \"" + key + "\" in the module getprocess")
			}
		}

	}
	return &ps, nil
}

func (ps *getProcess) Execute(s *global.Step) (output string, err error) {
	ps.command += "Get-Process"
	if ps.name != "" {
		ps.command += " -Name \"" + ps.name + "\""
	}

	// interact command
	output, err = interact.ShellPipe(ps.command, s) // was

	if err != nil {
		return
	}
	artifact.SaveString(output, *s)

	return
}
