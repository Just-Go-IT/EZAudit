package linux

import (
	"Just-Go-IT/EZAudit/artifact"
	"Just-Go-IT/EZAudit/global"
	"Just-Go-IT/EZAudit/interact"
	"Just-Go-IT/EZAudit/registry"
	"errors"
)

func init() {
	registry.Register("firewall", &firewall{}, false, registry.Linux)
}

type firewall struct {
	state          bool
	getDefaultZone bool
	getActiveZone  bool
	command        string
}

func (f firewall) New(p map[string]interface{}) (global.Module, error) {
	if len(p) >= 1 && len(p) <= 1 {
		return nil, errors.New("there is no action option in the module firewall")
	}

	ok := false

	// Checks for optional parameter Mode and check if the keys are allowed
	for key := range p {
		switch key {
		case "state":
			f.state, ok = p["state"].(bool)
			if !ok {
				return nil, errors.New("the key \"state\" is set for the module. The value must be a \"bool\"")
			}
		case "getDefaultZone":
			f.getDefaultZone, ok = p["getDefaultZone"].(bool)
			if !ok {
				return nil, errors.New("the key \"getDefaultZone\" is set for the module. The value must be a \"bool\"")
			}
		case "getActiveZone":
			f.getActiveZone, ok = p["getActiveZone"].(bool)
			if !ok {
				return nil, errors.New("the key \"getActiveZone\" is set for the module. The value must be a \"bool\"")
			}
		default:
			if key != "state" && key != "getDefaultZone" && key != "getActiveZone" {
				return nil, errors.New("there is no key called: \"" + key + "\" in the module firewall")
			}
		}
	}

	return &f, nil
}

func (f *firewall) Execute(s *global.Step) (output string, err error) {
	f.command = "firewall-cmd "

	if f.state {
		f.command += "--state"
	}
	if f.getActiveZone {
		f.command += "-get-active zone"
	}
	if f.getDefaultZone {
		f.command += "-get-default zone"
	}

	output, err = interact.ShellPipe(f.command, s)

	if err != nil {
		return
	}

	artifact.SaveString(output, *s)

	return
}
