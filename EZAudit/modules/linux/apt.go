package linux

import (
	"Just-Go-IT/EZAudit/artifact"
	"Just-Go-IT/EZAudit/global"
	"Just-Go-IT/EZAudit/interact"
	"Just-Go-IT/EZAudit/registry"
	"errors"
)

func init() {
	registry.Register("apt", &apt{}, false, registry.Linux)
}

type apt struct {
	cache       bool
	upgrade     bool
	purge       bool
	keyList     bool
	packageName string
	command     string
}

func (a apt) New(p map[string]interface{}) (global.Module, error) {
	ok := false

	// Checks for optional parameters and if the keys are supported
	for key, _ := range p {
		switch key {
		case "cache":
			a.cache, ok = p["cache"].(bool)
			if !ok {
				return nil, errors.New("the key \"cache\" is set for the module. The value must be a \"bool\"")
			}
		case "upgrade":
			a.upgrade, ok = p["upgrade"].(bool)
			if !ok {
				return nil, errors.New("the key \"upgrade\" is set for the module. The value must be a \"bool\"")
			}
		case "purge":
			a.purge, ok = p["purge"].(bool)
			if !ok {
				return nil, errors.New("the key \"purge\" is set for the module. The value must be a \"bool\"")
			}
		case "keyList":
			a.keyList, ok = p["keyList"].(bool)
			if !ok {
				return nil, errors.New("the key \"keyList\" is set for the module. The value must be a \"bool\"")
			}
		case "packageName":
			a.packageName, ok = p["packageName"].(string)
			if !ok {
				return nil, errors.New("the key \"packageName\" is set for the module. The value must be a \"string\"")
			}
		default:
			if key != "cache" && key != "upgrade" && key != "purge" && key != "keyList" && key != "packageName" {
				return nil, errors.New("there is no key called: \"" + key + "\" in the module apt")
			}
		}
	}
	if (!a.cache && !a.purge && !a.upgrade && !a.keyList && a.packageName != "") || (a.cache && a.purge && a.upgrade && a.keyList && a.packageName == "packageName") {
		return nil, errors.New("there is no action option in the module apt")
	}

	return &a, nil
}

func (a *apt) Execute(s *global.Step) (output string, err error) {
	a.command += "apt"

	if a.cache && a.packageName != "" {
		a.command += "-cache " + a.packageName
	}
	if a.purge && a.packageName != "" {
		a.command += " purge " + a.packageName
	}
	if a.upgrade {
		a.command += " -s upgrade"
	}
	if a.keyList {
		a.command += "-key list"
	}

	output, err = interact.ShellPipe(a.command, s)
	if err != nil {
		return
	}

	artifact.SaveString(output, *s)

	return
}
