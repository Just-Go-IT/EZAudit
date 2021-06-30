package linux

import (
	"Just-Go-IT/EZAudit/artifact"
	"Just-Go-IT/EZAudit/global"
	"Just-Go-IT/EZAudit/interact"
	"Just-Go-IT/EZAudit/registry"
	"errors"
)

func init() {
	registry.Register("modprobe", &modprobe{}, false, registry.Linux)
}

type modprobe struct {
	all        bool
	dryRun     bool
	version    bool
	verbose    bool
	moduleName string
	command    string
}

func (m modprobe) New(p map[string]interface{}) (global.Module, error) {
	ok := false

	m.moduleName, ok = p["moduleName"].(string)
	if !ok {
		return nil, errors.New("the key \"moduleName\" must be set and the value must be a \"string\"")
	}

	// Checks for optional parameters and if the keys are supported
	for key := range p {
		switch key {
		case "dryRun":
			m.dryRun, ok = p["dryRun"].(bool)
			if !ok {
				return nil, errors.New("the key \"dryRun\" is set for the module. The value must be a \"bool\"")
			}
		case "verbose":
			m.verbose, ok = p["verbose"].(bool)
			if !ok {
				return nil, errors.New("the key \"verbose\" is set for the module. The value must be a \"bool\"")
			}
		case "version":
			m.version, ok = p["version"].(bool)
			if !ok {
				return nil, errors.New("the key \"version\" is set for the module. The value must be a \"bool\"")
			}
		case "all":
			m.all, ok = p["all"].(bool)
			if !ok {
				return nil, errors.New("the key \"all\" is set for the module. The value must be a \"bool\"")
			}
		default:
			if key != "moduleName" && key != "dryRun" && key != "verbose" && key != "version" && key != "all" {
				return nil, errors.New("there is no key called: \"" + key + "\" in the module modprobe")
			}
		}
	}

	return &m, nil
}

func (m *modprobe) Execute(s *global.Step) (output string, err error) {
	m.command += "modprobe"
	if m.all {
		m.command += " -a"
	}
	if m.dryRun {
		m.command += " -n"
	}
	if m.verbose {
		m.command += " -V"
	}
	if m.version {
		m.command += " -v"
	}
	m.command += " " + m.moduleName

	output, err = interact.ShellPipe(m.command, s)
	if err != nil {
		return
	}

	artifact.SaveString(output, *s)
	return

}
