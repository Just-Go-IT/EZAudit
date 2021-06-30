package linux

import (
	"Just-Go-IT/EZAudit/artifact"
	"Just-Go-IT/EZAudit/global"
	"Just-Go-IT/EZAudit/interact"
	"Just-Go-IT/EZAudit/registry"
	"errors"
)

func init() {
	registry.Register("mount", &mount{}, false, registry.Linux)
}

type mount struct {
	admin   bool
	typ     string
	options string
	command string
}

func (m mount) New(p map[string]interface{}) (global.Module, error) {
	ok := false

	// Checks for optional parameters and if the keys are supported
	for key := range p {
		switch key {
		case "admin":
			m.admin, ok = p["admin"].(bool)
			if !ok {
				return nil, errors.New("the key \"admin\" is set for the module. The value must be a \"bool\"")
			}
		case "type":
			m.typ, ok = p["type"].(string)
			if !ok {
				return nil, errors.New("the key \"type\" is set for the module. The value must be a \"string\"")
			}
		case "options":
			m.options, ok = p["options"].(string)
			if !ok {
				return nil, errors.New("the key \"options\" is set for the module. The value must be a \"string\"")
			}
		default:
			if key != "admin" && key != "type" && key != "options" {
				return nil, errors.New("there is no key called: \"" + key + "\" in the module mount")
			}
		}
	}

	return &m, nil
}

func (m *mount) Execute(s *global.Step) (output string, err error) {
	if m.admin {
		m.command += "sudo"
	}
	m.command = " mount"
	if m.typ != "" {
		m.command += " -t " + m.typ
	}
	if m.options != "" {
		m.command += " -o " + m.options
	}

	output, err = interact.ShellPipe(m.command, s)
	if err != nil {
		return
	}

	artifact.SaveString(output, *s)
	return
}
