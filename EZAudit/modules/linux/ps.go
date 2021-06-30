package linux

import (
	"Just-Go-IT/EZAudit/artifact"
	"Just-Go-IT/EZAudit/global"
	"Just-Go-IT/EZAudit/interact"
	"Just-Go-IT/EZAudit/registry"
	"errors"
)

func init() {
	registry.Register("ps", &ps{}, false, registry.Linux)
}

type ps struct {
	addSecurityData bool
	standardSyntax  bool
	aux             bool
	target          string
	command         string
}

func (ps ps) New(p map[string]interface{}) (global.Module, error) {
	ok := false

	// Checks for optional parameters and if the keys are supported
	for key := range p {
		switch key {
		case "addSecurityData":
			ps.addSecurityData, ok = p["addSecurityData"].(bool)
			if !ok {
				return nil, errors.New("the key \"addSecurityData\" is set and the value must be a \"bool\"")
			}
		case "standardSyntax":
			ps.standardSyntax, ok = p["standardSyntax"].(bool)
			if !ok {
				return nil, errors.New("the key \"standardSyntax\" is set and the value must be a \"bool\"")
			}
		case "aux":
			ps.aux, ok = p["aux"].(bool)
			if !ok {
				return nil, errors.New("the key \"aux\" is set and the value must be a \"bool\"")
			}
		case "target":
			ps.target, ok = p["target"].(string)
			if !ok {
				return nil, errors.New("the key \"target\" is set and the value must be a \"string\"")
			}
		default:
			if key != "addSecurityData" && key != "standardSyntax" && key != "aux" && key != "target" {
				return nil, errors.New("there is no key called: \"" + key + "\" in the module ps")
			}
		}
	}
	if (!ps.addSecurityData && !ps.standardSyntax && !ps.aux) || (ps.standardSyntax && ps.addSecurityData && ps.aux) {
		return nil, errors.New("there is no action option in the module ps")
	}

	return &ps, nil
}

func (ps *ps) Execute(s *global.Step) (output string, err error) {
	if ps.target != "" {
		ps.command += ps.target
	}
	ps.command += "ps "
	if ps.addSecurityData {
		ps.command += "-eZ "
	}
	if ps.standardSyntax {
		ps.command += " -ef"
	}
	if ps.aux {
		ps.command += " aux"
	}

	output, err = interact.ShellPipe(ps.command, s)
	if err != nil {
		return
	}

	artifact.SaveString(output, *s)
	return
}
