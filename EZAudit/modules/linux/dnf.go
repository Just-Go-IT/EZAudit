package linux

import (
	"Just-Go-IT/EZAudit/artifact"
	"Just-Go-IT/EZAudit/global"
	"Just-Go-IT/EZAudit/interact"
	"Just-Go-IT/EZAudit/registry"
	"errors"
)

func init() {
	registry.Register("dnf", &dnf{}, false, registry.Linux)
}

type dnf struct {
	command     string
	repoList    bool
	checkUpdate bool
}

func (d dnf) New(p map[string]interface{}) (global.Module, error) {
	ok := false
	// Checks for optional parameters and if the keys are supported
	for key := range p {
		switch key {
		case "repoList":
			d.repoList, ok = p["repoList"].(bool)
			if !ok {
				return nil, errors.New("the key \"repoList\" is set for the module. The value must be a \"bool\"")
			}
		case "checkUpdate":
			d.checkUpdate, ok = p["checkUpdate"].(bool)
			if !ok {
				return nil, errors.New("the key \"checkUpdate\" is set for the module. The value must be a \"bool\"")
			}

		default:
			if key != "repoList" && key != "checkUpdate" {
				return nil, errors.New("there is no key called: \"" + key + "\" in the module dnf")
			}
		}
	}

	if (!d.repoList && !d.checkUpdate) || (d.repoList && d.checkUpdate) {
		return nil, errors.New("there is no action option in the module dpkg")
	}

	return &d, nil
}

func (d *dnf) Execute(s *global.Step) (output string, err error) {
	d.command = "dnf "

	if d.repoList {
		d.command += "repolist"
	}
	if d.checkUpdate {
		d.command += "check-update"
	}

	output, err = interact.ShellPipe(d.command, s)

	if err != nil {
		return
	}

	artifact.SaveString(output, *s)

	return
}
