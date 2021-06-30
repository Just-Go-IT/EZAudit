package linux

import (
	"Just-Go-IT/EZAudit/artifact"
	"Just-Go-IT/EZAudit/global"
	"Just-Go-IT/EZAudit/interact"
	"Just-Go-IT/EZAudit/registry"
	"errors"
)

func init() {
	registry.Register("touch", &touch{}, false, registry.Linux)
}

type touch struct {
	path    string
	command string
}

func (t touch) New(p map[string]interface{}) (global.Module, error) {
	ok := false
	t.path, ok = p["path"].(string)
	if !ok {
		return nil, errors.New("the key \"path\" must be set and the value must be a \"string\"")
	}

	// Checks for optional parameters and if the keys are supported
	for key := range p {
		if key != "path" {
			return nil, errors.New("there is no key called: \"" + key + "\" in the module ls")
		}
	}

	return &t, nil
}

func (t *touch) Execute(s *global.Step) (output string, err error) {
	t.command += "touch "

	if t.path != "" {
		t.command += t.path
	}

	output, err = interact.ShellPipe(t.command, s)

	if err != nil {
		return
	}

	artifact.SaveString(output, *s)
	return
}
