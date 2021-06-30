package linux

import (
	"Just-Go-IT/EZAudit/artifact"
	"Just-Go-IT/EZAudit/global"
	"Just-Go-IT/EZAudit/interact"
	"Just-Go-IT/EZAudit/registry"
	"errors"
)

func init() {
	registry.Register("crontab", &crontab{}, false, registry.Linux)
}

type crontab struct {
	user    string
	admin   bool
	list    bool
	editor  bool
	command string
}

func (c crontab) New(p map[string]interface{}) (global.Module, error) {
	ok := false

	c.user, ok = p["user"].(string)
	if !ok {
		return nil, errors.New("the key \"user\" must be set and the value must be a \"string\"")
	}

	c.list, ok = p["list"].(bool)
	if !ok {
		return nil, errors.New("the key \"list\" must be set and the value must be a \"bool\"")
	}

	// Checks for optional parameters and if the keys are supported
	for key := range p {
		switch key {
		case "admin":
			c.admin, ok = p["admin"].(bool)
			if !ok {
				return nil, errors.New("the key \"admin\" is set and the value must be a \"bool\"")
			}
		default:
			if key != "admin" && key != "list" && key != "user" {
				return nil, errors.New("there is no key called: \"" + key + "\" in the module crontab")
			}
		}
	}

	if !c.admin && c.user == "root" {
		return nil, errors.New("for User to be set as Root, we require Admin rights for this operation")
	}

	return &c, nil
}

func (c *crontab) Execute(s *global.Step) (output string, err error) {
	c.command = ""
	if c.admin {
		c.command += "sudo "
	}
	c.command += "crontab"
	if c.user != "" {
		c.command += " -u " + c.user
	}
	if c.editor {
		c.command += " -e"
	}
	if c.list {
		c.command += " -l"
	}

	output, err = interact.ShellPipe(c.command, s)
	if err != nil {
		return
	}

	artifact.SaveString(output, *s)

	return
}
