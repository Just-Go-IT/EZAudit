package linux

import (
	"Just-Go-IT/EZAudit/artifact"
	"Just-Go-IT/EZAudit/global"
	"Just-Go-IT/EZAudit/interact"
	"Just-Go-IT/EZAudit/registry"
	"errors"
)

func init() {
	registry.Register("whoami", &whoami{}, false, registry.Linux)
}

type whoami struct {
	version bool
	command string
}

func (w whoami) New(p map[string]interface{}) (global.Module, error) {
	ok := false

	// Checks for optional parameters and if the keys are supported
	for key := range p {
		switch key {
		case "version":
			w.version, ok = p["version"].(bool)
			if !ok {
				return nil, errors.New("the key \"version\" must be set and the value must be a \"bool\"")
			}
		default:
			if key != "version" {
				return nil, errors.New("there is no key called: \"" + key + "\" in the module")
			}
		}
	}

	return &w, nil
}

func (w *whoami) Execute(s *global.Step) (output string, err error) {
	w.command += "whoami"

	if w.version {
		w.command += " --version"
	}

	output, err = interact.ShellPipe(w.command, s)
	if err != nil {
		return
	}

	artifact.SaveString(output, *s)

	return
}
