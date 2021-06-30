package linux

import (
	"Just-Go-IT/EZAudit/artifact"
	"Just-Go-IT/EZAudit/global"
	"Just-Go-IT/EZAudit/interact"
	"Just-Go-IT/EZAudit/registry"
	"errors"
)

func init() {
	registry.Register("grpck", &grpck{}, true, registry.Linux)
}

type grpck struct {
	admin        bool
	readOnlyMode bool
	optionQ      bool
	group        string
	shadow       string
	command      string
}

func (g grpck) New(p map[string]interface{}) (global.Module, error) {
	ok := false

	// Checks for optional parameters and if the keys are supported
	for key := range p {
		switch key {
		case "admin":
			g.admin, ok = p["admin"].(bool)
			if !ok {
				return nil, errors.New("the key \"admin\" is set for the module. The value must be a \"bool\"")
			}
		case "readOnlyMode":
			g.readOnlyMode, ok = p["readOnlyMode"].(bool)
			if !ok {
				return nil, errors.New("the key \"readOnlyMode\" is set for the module. The value must be a \"bool\"")
			}
		case "optionQ":
			g.optionQ, ok = p["optionQ"].(bool)
			if !ok {
				return nil, errors.New("the key \"optionQ\" is set for the module. The value must be a \"bool\"")
			}
		case "group":
			g.group, ok = p["group"].(string)
			if !ok {
				return nil, errors.New("the key \"group\" is set for the module. The value must be a \"string\"")
			}
		case "shadow":
			g.shadow, ok = p["shadow"].(string)
			if !ok {
				return nil, errors.New("the key \"shadow\" is set for the module. The value must be a \"string\"")
			}
		default:
			if key != "admin" && key != "readonlymode" && key != "optionQ" && key != "group" && key != "shadow" {
				return nil, errors.New("there is no key called: \"" + key + "\" in the module dpkg")
			}
		}
	}

	return &g, nil
}

func (g *grpck) Execute(s *global.Step) (output string, err error) {
	if g.admin {
		g.command += "sudo "
	}
	if g.group != "" {
		g.command += " " + g.group
	}
	if g.shadow != "" {
		g.command += " " + g.shadow
	}
	g.command += "grpck"
	if g.optionQ {
		g.command += " -q"
	}
	if g.readOnlyMode {
		g.command += " -r"
	}

	output, err = interact.ShellPipe(g.command, s)
	if err != nil {
		return
	}

	artifact.SaveString(output, *s)

	return
}
