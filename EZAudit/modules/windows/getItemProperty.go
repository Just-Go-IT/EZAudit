package windows

import (
	"Just-Go-IT/EZAudit/global"
	"Just-Go-IT/EZAudit/interact"
	"Just-Go-IT/EZAudit/registry"
	"errors"
)

func init() {
	registry.Register("getItemProperty", &getItemProperty{}, false, registry.Windows)
}

type getItemProperty struct {
	path        string
	literalPath string
	name        string
	filter      string
	include     string
	exclude     string
	command     string
}

func (c getItemProperty) New(p map[string]interface{}) (global.Module, error) {
	// Parse arguments
	ok := false
	c.path, ok = p["path"].(string)
	if !ok {
		c.literalPath, ok = p["literalPath"].(string)
		if !ok {
			return nil, errors.New("the parameter \"path\" or \"literalPath\" must be set and the value must be a \"string\"")
		}
	}
	// Check for optional parameter Mode an look
	for k, _ := range p {
		switch k {
		case "name":
			c.name, ok = p["name"].(string)
			if !ok {
				return nil, errors.New("the key \"name\" is set but the value, which should be a \"string\" could not be parsed")
			}
		case "filter":
			c.filter, ok = p["filter"].(string)
			if !ok {
				return nil, errors.New("the key \"filter\" is set but the value, which should be a \"string\" could not be parsed")
			}
		case "include":
			c.include, ok = p["include"].(string)
			if !ok {
				return nil, errors.New("the key \"include\" is set but the value, which should be a \"string\" could not be parsed")
			}
		case "exclude":
			c.exclude = p["exclude"].(string)
			if !ok {
				return nil, errors.New("the key \"exclude\" is set but the value, which should be a \"string\" could not be parsed")
			}
		default:
			if k != "path" && k != "literalPath" && k != "name" && k != "filter" && k != "include" && k != "exclude" {
				return nil, errors.New("there is no key called \"" + k + "\" in the module getItemProperty")
			}
		}

	}
	return &c, nil
}

func (c *getItemProperty) Execute(s *global.Step) (output string, err error) {
	// build command
	c.command += "Get-ItemProperty"
	if c.path != "" {
		c.command += " -Path \"" + c.path + "\""
	} else {
		c.command += " -LiteralPath " + "\"" + c.literalPath + "\""
	}
	if c.name != "" {
		c.command += " -Name " + "\"" + c.name + "\""
	}
	if c.filter != "" {
		c.command += " -Filter " + c.filter
	}
	if c.include != "" {
		c.command += " -Include " + c.include
	}
	if c.exclude != "" {
		c.command += " -Exclude " + c.exclude
	}

	// interact command
	output, err = interact.ShellPipe(c.command, s)

	return
}
