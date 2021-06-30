package linux

import (
	"Just-Go-IT/EZAudit/artifact"
	"Just-Go-IT/EZAudit/global"
	"Just-Go-IT/EZAudit/interact"
	"Just-Go-IT/EZAudit/registry"
	"errors"
)

func init() {
	registry.Register("iptables", &iptables{}, true, registry.Linux)
}

type iptables struct {
	admin    bool
	verbose  bool
	nummeric bool
	options  string
	command  string
}

func (i iptables) New(p map[string]interface{}) (global.Module, error) {
	ok := false

	i.admin, ok = p["admin"].(bool)
	if !ok {
		return nil, errors.New("the key \"admin\" must be set and the value must be a \"bool\"")
	}

	// Checks for optional parameters and if the keys are supported
	for key := range p {
		switch key {
		case "options":
			i.options, ok = p["options"].(string)
			if !ok {
				return nil, errors.New("the key \"options\" is set for the module. The value must be a \"string\"")
			}
		case "verbose":
			i.verbose, ok = p["verbose"].(bool)
			if !ok {
				return nil, errors.New("the key \"verbose\" is set for the module. The value must be a \"bool\"")
			}
		case "nummeric":
			i.nummeric, ok = p["nummeric"].(bool)
			if !ok {
				return nil, errors.New("the key \"nummeric\" is set for the module. The value must be a \"bool\"")
			}
		default:
			if key != "options" && key != "verbose" && key != "nummeric" && key != "admin" {
				return nil, errors.New("there is no key called: \"" + key + "\" in the module iptables")
			}
		}
	}
	if !i.admin {
		return nil, errors.New("you need to be root to execute iptables")
	}

	return &i, nil
}

func (i *iptables) Execute(s *global.Step) (output string, err error) {
	i.command += "sudo iptables -L "

	if i.options != "" {
		i.command += i.options + " "
	}
	if i.verbose {
		i.command += "-v"
	}
	if i.nummeric {
		i.command += "-n"
	}

	output, err = interact.ShellPipe(i.command, s)
	if err != nil {
		return
	}

	artifact.SaveString(output, *s)

	return
}
