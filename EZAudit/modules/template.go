package modules

import (
	"Just-Go-IT/EZAudit/artifact"
	"Just-Go-IT/EZAudit/global"
	"Just-Go-IT/EZAudit/interact"
	"Just-Go-IT/EZAudit/registry"
	"errors"
)

func init() {
	registry.Register("moduleName", &moduleTemplate{}, false, registry.Windows|registry.Linux|registry.SupportAll)
}

type moduleTemplate struct {
	argument1 bool
	argument2 string
	argument3 int
	command   string
}

func (t moduleTemplate) New(p map[string]interface{}) (global.Module, error) {
	ok := false

	// Needs a parameter
	t.argument2, ok = p["needed"].(string)
	if !ok {
		return nil, errors.New("the key \"needed\" must be set and the value must be a \"type\"")
	}

	// Checks for optional parameters and if the keys are supported
	for key := range p {
		switch key {
		case "path":
			t.argument1, ok = p["argument1"].(bool)
			if !ok {
				return nil, errors.New("the key \"argument1\" must be set and the value must be a \"bool\"")
			}
		default:
			if key != "argument1" {
				return nil, errors.New("there is no key called: \"" + key + "\" in the module")
			}
		}
	}
	return &t, nil
}

func (t *moduleTemplate) Execute(s *global.Step) (output string, err error) {
	t.command = ""

	if t.argument1 {
		t.command += " "
	}

	output, err = interact.ShellPipe(t.command, s)
	if err != nil {
		return
	}

	artifact.SaveString(output, *s)

	return
}
