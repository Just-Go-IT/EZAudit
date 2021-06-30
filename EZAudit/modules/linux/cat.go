package linux

import (
	"Just-Go-IT/EZAudit/artifact"
	"Just-Go-IT/EZAudit/global"
	"Just-Go-IT/EZAudit/interact"
	"Just-Go-IT/EZAudit/registry"
	"errors"
)

func init() {
	registry.Register("cat", &cat{}, false, registry.Linux)
}

type cat struct {
	path    string
	command string
}

func (cat cat) New(p map[string]interface{}) (global.Module, error) {
	var ok bool

	// Checks for optional parameters and if the keys are supported
	for key := range p {
		switch key {
		case "path":
			cat.path, ok = p["path"].(string)
			if !ok {
				return nil, errors.New("the key \"path\" must be set and the value must be a\"string\"")
			}
		default:
			if key != "path" {
				return nil, errors.New("there is no key called: \"" + key + "\" in the module getContent")
			}
		}
	}

	return &cat, nil
}

func (cat *cat) Execute(s *global.Step) (output string, err error) {
	cat.command = "cat " + cat.path
	output, err = interact.ShellPipe(cat.command, s)

	if err != nil {
		return
	}

	artifact.SaveFile(cat.path, *s)

	return
}
