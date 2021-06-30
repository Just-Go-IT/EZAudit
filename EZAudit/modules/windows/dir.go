package windows

import (
	"Just-Go-IT/EZAudit/artifact"
	"Just-Go-IT/EZAudit/global"
	"Just-Go-IT/EZAudit/interact"
	"Just-Go-IT/EZAudit/registry"
	"errors"
)

func init() {
	registry.Register("dir", &dir{}, false, registry.Windows)
}

type dir struct {
	path    string
	command string
}

func (d dir) New(p map[string]interface{}) (global.Module, error) {
	// Parse arguments
	ok := true
	// Check for optional parameter Mode and whitelist the keys
	for key := range p {
		switch key {
		case "path":
			d.path, ok = p["path"].(string)
			if !ok {
				return nil, errors.New("the key \"path\" is set for the module. The value must be a \"string\"")
			}
		default:
			if key != "path" {
				return nil, errors.New("there is no key called: \"" + key + "\" in the module dir")
			}
		}

	}
	return &d, nil
}

func (d *dir) Execute(s *global.Step) (output string, err error) {
	d.command += "dir"
	if d.path != "" {
		d.command += " -Path \"" + d.path + "\""
	}

	// interact command
	output, err = interact.ShellPipe(d.command, s)
	if err != nil {
		return
	}

	artifact.SaveString(output, *s)

	return
}
