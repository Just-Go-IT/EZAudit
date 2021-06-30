package windows

import (
	"Just-Go-IT/EZAudit/artifact"
	"Just-Go-IT/EZAudit/global"
	"Just-Go-IT/EZAudit/interact"
	"Just-Go-IT/EZAudit/registry"
	"errors"
)

func init() {
	registry.Register("getContent", &getContent{}, false, registry.Windows)
}

type getContent struct {
	path    string
	command string
}

func (gc getContent) New(p map[string]interface{}) (global.Module, error) {
	// Parse arguments
	ok := true
	for key, _ := range p {
		switch key {
		case "path":
			gc.path, ok = p["path"].(string)
			if !ok {
				return nil, errors.New("the key \"path\" must be set and the value must be a\"string\"")
			}
		default:
			if key != "path" {
				return nil, errors.New("there is no key called: \"" + key + "\" in the module getContent")
			}
		}
	}

	return &gc, nil
}

func (gc *getContent) Execute(s *global.Step) (output string, err error) {
	gc.command += "Get-Content"

	if gc.path != "" {
		gc.command += " -Path \"" + gc.path + "\""
	}

	// interact command
	output, err = interact.ShellPipe(gc.command, s)

	if err != nil {
		return
	}

	artifact.SaveString(output, *s)

	return
}
