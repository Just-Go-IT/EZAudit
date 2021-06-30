package linux

import (
	"Just-Go-IT/EZAudit/artifact"
	"Just-Go-IT/EZAudit/global"
	"Just-Go-IT/EZAudit/interact"
	"Just-Go-IT/EZAudit/registry"
	"errors"
)

func init() {
	registry.Register("ip", &ip{}, false, registry.Linux)
}

type ip struct {
	command     string
	showAddress bool
	routeList   bool
	old         string
}

func (i ip) New(p map[string]interface{}) (global.Module, error) {
	ok := false

	// Checks for optional parameters and if the keys are supported
	for key := range p {
		switch key {
		case "showAddress":
			i.showAddress, ok = p["showAddress"].(bool)
			if !ok {
				return nil, errors.New("the key \"showAddress\" is set and the value must be a \"bool\"")
			}
		case "routeList":
			i.routeList, ok = p["routeList"].(bool)
			if !ok {
				return nil, errors.New("the key \"routeList\" is set and the value must be a \"bool\"")
			}
		case "old":
			i.old, ok = p["old"].(string)
			if !ok {
				return nil, errors.New("the key \"old\" is set and the value must be a \"bool\"")
			}
		default:
			if key != "showAddress" && key != "routeList" && key != "old" {
				return nil, errors.New("there is no key called: \"" + key + "\" in the module ip")
			}
		}
	}

	if (!i.routeList && !i.showAddress && i.old == "") || (i.routeList && i.showAddress && i.old != "") {
		return nil, errors.New("there is no action option in the module ip")
	}

	return &i, nil
}
func (i *ip) Execute(currentStep *global.Step) (output string, err error) {
	i.command += "ip "
	if i.showAddress {
		i.command += "addr show"
	}
	if i.routeList {
		i.command += "route list"
	}
	if i.old != "" {
		i.command = i.old
	}

	output, err = interact.ShellPipe(i.command, currentStep)
	if err != nil {
		return
	}

	artifact.SaveString(output, *currentStep)

	return
}
