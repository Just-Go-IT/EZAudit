package linux

import (
	"Just-Go-IT/EZAudit/artifact"
	"Just-Go-IT/EZAudit/global"
	"Just-Go-IT/EZAudit/interact"
	"Just-Go-IT/EZAudit/registry"
	"errors"
)

func init() {
	registry.Register("ls", &ls{}, false, registry.Linux)
}

type ls struct {
	path    string
	command string
}

func (l ls) New(p map[string]interface{}) (global.Module, error) {
	ok := false
	// Checks for optional parameters and if the keys are supported
	for key := range p {
		switch key {
		case "path":
			l.path, ok = p["path"].(string)
			if !ok {
				return nil, errors.New("the key \"path\" must be set and the value must be a \"string\"")
			}
		default:
			if key != "path" {
				return nil, errors.New("there is no key called: \"" + key + "\" in the module ls")
			}
		}
	}

	return &l, nil
}

func (l *ls) Execute(s *global.Step) (output string, err error) {
	l.command += "ls -l "

	if l.path != "" {
		l.command += l.path
	}

	output, err = interact.ShellPipe(l.command, s)
	if err != nil {
		return
	}

	artifact.SaveString(output, *s)

	return
}
