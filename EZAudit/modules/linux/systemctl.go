package linux

import (
	"Just-Go-IT/EZAudit/artifact"
	"Just-Go-IT/EZAudit/global"
	"Just-Go-IT/EZAudit/interact"
	"Just-Go-IT/EZAudit/registry"
	"errors"
)

func init() {
	registry.Register("systemctl", &systemctl{}, false, registry.Linux)
}

type systemctl struct {
	now         bool
	disableUnit string
	command     string
	isEnabled   string
	status      string
}

func (s systemctl) New(p map[string]interface{}) (global.Module, error) {
	ok := false

	// Checks for optional parameters and if the keys are supported
	for key := range p {
		switch key {
		case "now":
			s.now, ok = p["now"].(bool)
			if !ok {
				return nil, errors.New("the key \"now\" is set for the module. The value must be a \"bool\"")
			}
		case "disableUnit":
			s.disableUnit, ok = p["disableUnit"].(string)
			if !ok {
				return nil, errors.New("the key \"disableUnit\" is set for the module. The value must be a \"string\"")
			}
		case "isEnabled":
			s.isEnabled, ok = p["isEnabled"].(string)
			if !ok {
				return nil, errors.New("the key \"isEnabled\" is set for the module. The value must be a \"string\"")
			}
		case "status":
			s.status, ok = p["status"].(string)
			if !ok {
				return nil, errors.New("the key \"status\" is set for the module. The value must be a \"string\"")
			}
		default:
			if key != "disableUnit" && key != "now" && key != "isEnabled" && key != "status" {
				return nil, errors.New("there is no key called: \"" + key + "\" in the module systemctl")
			}
		}
	}

	return &s, nil
}

func (s *systemctl) Execute(step *global.Step) (output string, err error) {
	s.command += "systemctl"

	if s.now {
		s.command += " --now"
	}
	if s.disableUnit != "" {
		s.command += " disable " + s.disableUnit
	}
	if s.isEnabled != "" {
		s.command += " is-enabled " + s.isEnabled
	}
	if s.status != "" {
		s.command += " status " + s.status
	}

	output, err = interact.ShellPipe(s.command, step)

	if err != nil && err.Error() != "exit status 1" {
		return
	} else if err != nil && err.Error() == "exit status 1" {
		err = nil
	}

	artifact.SaveString(output, *step)
	return
}
