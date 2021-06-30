package windows

import (
	"Just-Go-IT/EZAudit/global"
	"Just-Go-IT/EZAudit/interact"
	"Just-Go-IT/EZAudit/registry"
	"errors"
)

func init() {
	registry.Register("getGpo", &getGPO{}, false, registry.Windows)
}

type getGPO struct {
	all     bool   //Indicates that the cmdlet gets all the GPOs in the domain.
	domain  string //Specifies the domain for this cmdlet. You must specify the fully qualified domain name (FQDN) of the domain.
	guid    string //Specifies the GPO to retrieve by its globally unique identifier (GUID). The GUID uniquely identifies the GPO.
	name    string //Specifies the GPO to retrieve by its display name.
	server  string //Specifies the name of the domain controller that this cmdlet contacts to complete the operation. You can specify either the fully qualified domain name (FQDN) or the host name.
	command string
}

func (c getGPO) New(p map[string]interface{}) (global.Module, error) {
	// Parse arguments
	ok := false
	c.guid, ok = p["guid"].(string)
	if !ok {
		c.name, ok = p["name"].(string)
		if !ok {
			c.domain, ok = p["domain"].(string)
			if !ok {
				c.server, ok = p["server"].(string)
				if !ok {
					return nil, errors.New(" GetGPO needs at least one of the the parameters: \"guid\", \"name\", \"domain\", \"server\". Required Type:\"string\"")
				}
			}
		}
	}
	// Check for optional parameter
	for k, _ := range p {
		switch k {
		case "guid":
			c.guid, ok = p["guid"].(string)
			if !ok {
				return nil, errors.New("the key \"guid\" is set but the value, which should be a \"string\" could not be parsed")
			}
		case "name":
			c.name, ok = p["name"].(string)
			if !ok {
				return nil, errors.New("the key \"name\" is set but the value, which should be a \"string\" could not be parsed")
			}
		case "domain":
			c.domain, ok = p["domain"].(string)
			if !ok {
				return nil, errors.New("the key \"domain\" is set but the value, which should be a \"string\" could not be parsed")
			}
		case "server":
			c.server = p["server"].(string)
			if !ok {
				return nil, errors.New("the key \"server\" is set but the value, which should be a \"string\" could not be parsed")
			}
		case "all":
			c.all = p["all"].(bool)
			if !ok {
				return nil, errors.New("the key \"all\" is set but the value, which should be a \"bool\" could not be parsed")
			}
		default:
			return nil, errors.New("there is no key called \"" + k + "\" in the module getItemProperty")
		}
	}
	return &c, nil
}

func (c *getGPO) Execute(s *global.Step) (output string, err error) {
	// build command
	c.command += "Get-GPO"
	if c.guid != "" {
		c.command += " -Guid " + c.guid
	} else if c.name != "" {
		c.command += " -Name " + "\"" + c.name + "\""
	}
	if c.domain != "" {
		c.command += " -Domain " + "\"" + c.domain + "\""
	}
	if c.server != "" {
		c.command += " -Server " + "\"" + c.server + "\""
	}

	if c.all {
		c.command += " -All"
	}

	// interact command
	output, err = interact.ShellPipe(c.command, s)

	return
}
