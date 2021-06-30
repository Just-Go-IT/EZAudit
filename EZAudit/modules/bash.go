// Package modules includes all modules which are usable for every os
package modules

import (
	"Just-Go-IT/EZAudit/global"
	"Just-Go-IT/EZAudit/interact"
	"Just-Go-IT/EZAudit/registry"
	"errors"
)

func init() {
	registry.Register("bash", &bash{}, false, registry.Linux|registry.Windows)
}

type bash struct {
	script  string
	path    string
	command string
}

func (b bash) New(p map[string]interface{}) (global.Module, error) {
	ok := false

	// Checks for optional parameters and if the keys are supported
	for key := range p {
		switch key {
		case "script":
			b.script, ok = p["script"].(string)
			if !ok {
				return nil, errors.New("the key \"script\" is set for the module. The value must be a \"string\"")
			}
		case "path":
			b.path, ok = p["path"].(string)
			if !ok {
				return nil, errors.New("the key \"path\" is set for the module. The value must be a \"string\"")
			}
		default:
			if key != "script" && key != "path" {
				return nil, errors.New("there is no key called: \"" + key + "\".")
			}
		}
	}

	if b.script == "" && b.path == "" {
		return nil, errors.New("there is no script and no path selected")
	}
	return &b, nil
}

func (b *bash) Execute(s *global.Step) (output string, err error) {

	if b.path != "" {
		b.command = "." + b.path
	} else {
		b.command = b.script
	}

	output, err = interact.ShellPipe(b.command, s)

	return
}
