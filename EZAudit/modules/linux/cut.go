package linux

import (
	"Just-Go-IT/EZAudit/artifact"
	"Just-Go-IT/EZAudit/global"
	"Just-Go-IT/EZAudit/interact"
	"Just-Go-IT/EZAudit/registry"
	"errors"
)

func init() {
	registry.Register("cut", &cut{}, false, registry.Linux)
}

type cut struct {
	characters bool
	delimiter  bool
	target     string
	command    string
}

func (c cut) New(p map[string]interface{}) (global.Module, error) {
	ok := false

	// Checks for optional parameters and if the keys are supported
	for key := range p {
		switch key {
		case "characters":
			c.characters, ok = p["characters"].(bool)
			if !ok {
				return nil, errors.New("the key \"characters\" is set for the module. The value must be a \"bool\"")
			}
		case "delimiter":
			c.delimiter, ok = p["delimiter"].(bool)
			if !ok {
				return nil, errors.New("the key \"delimiter\" is set for the module. The value must be a \"bool\"")
			}
		case "target":
			c.target, ok = p["target"].(string)
			if !ok {
				return nil, errors.New("the key \"target\" is set for the module. The value must be a \"string\"")
			}
		default:
			if key != "characters" && key != "delimiter" && key != "target" && key != "usePipe" {
				return nil, errors.New("there is no key called: \"" + key + "\"")
			}
		}
	}
	if c.target == "" {
		return nil, errors.New("no target used")
	}

	if !c.delimiter && !c.characters && c.target == "" {
		return nil, errors.New("there is no action option in the module cut")
	}

	return &c, nil
}

func (c *cut) Execute(s *global.Step) (output string, err error) {
	c.command += "cut "

	if c.characters {
		c.command += "-c "
	}
	if c.delimiter {
		c.command += "-d "
	}
	if c.target != "" {
		c.command += c.target
	}

	output, err = interact.ShellPipe(c.command, s)
	if err != nil {
		return
	}

	artifact.SaveFile(c.target, *s)
	artifact.SaveString(output, *s)

	return
}
