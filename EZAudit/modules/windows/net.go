package windows

import (
	"Just-Go-IT/EZAudit/artifact"
	"Just-Go-IT/EZAudit/global"
	"Just-Go-IT/EZAudit/interact"
	"Just-Go-IT/EZAudit/registry"
	"errors"
)

func init() {
	registry.Register("net", &net{}, false, registry.Windows)
}

type net struct {
	accounts   bool // lists password and logon requirements for users
	view       bool // lists all the networked devices
	group      bool // lists groups (can only be used on domain controller)
	localGroup bool //lists local groups
	start      bool // lists all started Windows services
	share      bool // lists shared resources
	user       bool // lists local users
	command    string
}

func (c net) New(p map[string]interface{}) (global.Module, error) {
	// Parse arguments
	ok := false
	i := 0
	// Check for optional parameter
	for k := range p {
		switch k {
		case "accounts":
			c.accounts, ok = p["accounts"].(bool)
			if !ok {
				i++
				return nil, errors.New("the key \"accounts\" is set but the value, which should be a \"bool\" could not be parsed")
			}
		case "view":
			c.view, ok = p["view"].(bool)
			if !ok {
				i++
				return nil, errors.New("the key \"view\" is set but the value, which should be a \"bool\" could not be parsed")
			}
		case "group":
			c.group, ok = p["group"].(bool)
			if !ok {
				i++
				return nil, errors.New("the key \"group\" is set but the value, which should be a \"bool\" could not be parsed")
			}
		case "localGroup":
			c.localGroup = p["localGroup"].(bool)
			if !ok {
				i++
				return nil, errors.New("the key \"localGroup\" is set but the value, which should be a \"bool\" could not be parsed")
			}
		case "start":
			c.start = p["start"].(bool)
			if !ok {
				i++
				return nil, errors.New("the key \"start\" is set but the value, which should be a \"bool\" could not be parsed")
			}
		case "share":
			c.share = p["share"].(bool)
			if !ok {
				i++
				return nil, errors.New("the key \"share\" is set but the value, which should be a \"bool\" could not be parsed")
			}
		case "user":
			c.user = p["user"].(bool)
			if !ok {
				i++
				return nil, errors.New("the key \"user\" is set but the value, which should be a \"bool\" could not be parsed")
			}
		default:
			return nil, errors.New("there is no key called \"" + k + "\" in the module getItemProperty")
		}
	}
	if i > 1 {
		return nil, errors.New("the module supports only one \"key\" at a time")
	}
	return &c, nil
}

func (c *net) Execute(s *global.Step) (output string, err error) {
	// build command
	if c.accounts {
	} else if c.view {
		c.command = "net view"
	} else if c.group {
		c.command = "net group"
	} else if c.localGroup {
		c.command = "net localGroup"
	} else if c.start {
		c.command = "net start"
	} else if c.share {
		c.command = "net share"
	} else if c.user {
		c.command = "net user"
	}

	output, err = interact.ShellPipe(c.command, s)

	artifact.SaveString(output, *s)

	return
}
