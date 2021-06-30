package linux

import (
	"Just-Go-IT/EZAudit/artifact"
	"Just-Go-IT/EZAudit/global"
	"Just-Go-IT/EZAudit/interact"
	"Just-Go-IT/EZAudit/registry"
	"errors"
)

func init() {
	registry.Register("command", &command{}, false, registry.Linux)
}

type command struct {
	printDescription   bool //v
	verboseDescription bool //V
	target             string
	command            string
}

func (c command) New(p map[string]interface{}) (global.Module, error) {
	ok := false
	// Checks for optional parameters and if the keys are supported
	for key := range p {
		switch key {
		case "printDescription":
			c.printDescription, ok = p["printDescription"].(bool)
			if !ok {
				return nil, errors.New("the key \"printDiscription\" is set for the module. The value must be a \"bool\"")
			}
		case "verboseDescription":
			c.verboseDescription, ok = p["verboseDescription"].(bool)
			if !ok {
				return nil, errors.New("the key \"verboseDiscription\" is set for the module. The value must be a \"bool\"")
			}
		case "target":
			c.target, ok = p["target"].(string)
			if !ok {
				return nil, errors.New("the key \"target\" is set for the module. The value must be a \"string\"")
			}
		default:
			if key != "printDescription" && key != "verboseDescription" && key != "target" {
				return nil, errors.New("there is no key called: \"" + key + "\" in the module command")
			}
		}
	}

	if (!c.printDescription && !c.verboseDescription) || (c.printDescription && c.verboseDescription) {
		return nil, errors.New("there is no action option in the module command")
	}
	return &c, nil
}

func (c *command) Execute(s *global.Step) (output string, err error) {
	c.command = "command"

	if c.printDescription {
		c.command += " -v"
	}
	if c.verboseDescription {
		c.command += " -V"
	}
	if c.target != "" {
		c.command += " " + c.target
	}

	output, err = interact.ShellPipe(c.command, s)
	if err != nil {
		return
	}

	artifact.SaveString(output, *s)

	return
}
