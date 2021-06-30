package windows

import (
	"Just-Go-IT/EZAudit/global"
	"Just-Go-IT/EZAudit/interact"
	"Just-Go-IT/EZAudit/registry"
	"errors"
)

func init() {
	registry.Register("getItemPropertyValue", &getItemPropertyValue{}, false, registry.Windows)
}

type getItemPropertyValue struct {
	path        string
	literalPath string
	name        string
	filter      string
	include     string
	exclude     string
	command     string
}

func (c getItemPropertyValue) New(p map[string]interface{}) (global.Module, error) {
	// Parse arguments
	ok := false
	c.path, ok = p["path"].(string)
	if !ok {
		c.literalPath, ok = p["literalPath"].(string)
		if !ok {
			return nil, errors.New("the parameter \"path\" or \"literalPath\" must be set and the value must be a \"string\"")
		}
	}
	// Check for optional parameter Mode
	for k, _ := range p {
		switch k {
		case "path":
		case "literalPath":
		case "name":
			c.name, ok = p["name"].(string)
			if !ok {
				return nil, errors.New("the parameter \"" + k + "\" must be a string.")
			}
		case "filter":
			c.filter, ok = p["filter"].(string)
			if !ok {
				return nil, errors.New("the parameter \"" + k + "\" must be a string")
			}
		case "include":
			c.include, ok = p["include"].(string)
			if !ok {
				return nil, errors.New("the parameter \"" + k + "\" must be a string")
			}
		case "exclude":
			c.exclude = p["exclude"].(string)
			if !ok {
				return nil, errors.New("the parameter \"" + k + "\" must be a string")
			}
		default:
			return nil, errors.New("the parameter \"" + k + "\" is not supported")
		}

	}
	return &c, nil
}

func (c *getItemPropertyValue) Execute(s *global.Step) (output string, err error) {
	// build command
	c.command += "Get-ItemPropertyValue"
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
