package linux

import (
	"Just-Go-IT/EZAudit/artifact"
	"Just-Go-IT/EZAudit/global"
	"Just-Go-IT/EZAudit/interact"
	"Just-Go-IT/EZAudit/registry"
	"errors"
)

func init() {
	registry.Register("dpkg", &dpkg{}, false, registry.Linux)
}

type dpkg struct {
	search      bool
	status      bool
	verify      bool
	pack        string
	listPattern string
	command     string
}

func (d dpkg) New(p map[string]interface{}) (global.Module, error) {
	ok := false
	d.pack, ok = p["package"].(string)
	if !ok {
		return nil, errors.New("the key \"package\" must be set and the value must be a \"string\"")
	}

	// Checks for optional parameters and if the keys are supported
	for key := range p {
		switch key {
		case "search":
			d.search, ok = p["search"].(bool)
			if !ok {
				return nil, errors.New("the key \"search\" is set for the module. The value must be a \"bool\"")
			}
		case "status":
			d.status, ok = p["status"].(bool)
			if !ok {
				return nil, errors.New("the key \"status\" is set for the module. The value must be a \"bool\"")
			}
		case "verify":
			d.verify, ok = p["verify"].(bool)
			if !ok {
				return nil, errors.New("the key \"verify\" is set for the module. The value must be a \"bool\"")
			}
		case "listPattern":
			d.listPattern, ok = p["listPattern"].(string)
			if !ok {
				return nil, errors.New("the key \"listPattern\" is set for the module. The value must be a \"string\"")
			}
		default:
			if key != "package" && key != "search" && key != "status" && key != "verify" && key != "listPattern" {
				return nil, errors.New("there is no key called: \"" + key + "\" in the module dpkg")
			}
		}
	}

	if !d.verify && !d.search && !d.status && d.listPattern == "" {
		return nil, errors.New("there is no action option in the module dpkg")
	}

	return &d, nil
}

func (d *dpkg) Execute(currentStep *global.Step) (output string, err error) {
	d.command = "dpkg"

	if d.search {
		d.command += " -S "
	}
	if d.status {
		d.command += " -s "
	}
	if d.listPattern != "" {
		d.command += " -l " + d.listPattern
	}
	if d.verify {
		d.command += " --verify "
	}
	if d.pack != "" {
		d.command += d.pack
	}

	output, err = interact.ShellPipe(d.command, currentStep)

	if err != nil {
		return
	}

	artifact.SaveString(output, *currentStep)

	return
}
