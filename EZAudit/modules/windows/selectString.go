package windows

import (
	"Just-Go-IT/EZAudit/artifact"
	"Just-Go-IT/EZAudit/global"
	"Just-Go-IT/EZAudit/interact"
	"Just-Go-IT/EZAudit/registry"
	"errors"
)

func init() {
	registry.Register("selectstring", &selectString{}, false, registry.Windows)
}

type selectString struct {
	path          string
	pattern       string
	caseSensitive bool
	quiet         bool
	command       string
}

func (f selectString) New(p map[string]interface{}) (global.Module, error) {
	// Parse arguments
	ok := false

	f.pattern, ok = p["pattern"].(string)
	if !ok {
		return nil, errors.New("the key \"pattern\" must be set and the value must be a \"string\"")
	}

	// Check for optional parameter Mode and check if the keys are allowed
	for key, _ := range p {
		switch key {
		case "path":
			f.path, ok = p["path"].(string)
			if !ok {
				return nil, errors.New("the key \"path\" must be set and the value must be a \"string\"")
			}
		case "caseSensitive":
			f.caseSensitive, ok = p["caseSensitive"].(bool)
			if !ok {
				return nil, errors.New("the key \"caseSensitive\" is set for the module. The value must be a \"bool\"")
			}
		case "quiet":
			f.quiet, ok = p["quiet"].(bool)
			if !ok {
				return nil, errors.New("the key \"quiet\" is set for the module. The value must be a \"bool\"")
			}
		default:
			if key != "path" && key != "pattern" && key != "caseSensitive" && key != "quiet" {
				return nil, errors.New("the key called: \"" + key + "\" isn't supported")
			}
		}

	}
	return &f, nil
}

func (f *selectString) Execute(s *global.Step) (output string, err error) {
	f.command += "Select-String"
	if f.path != "" {
		f.command += " -Path \"" + f.path + "\""
	}
	f.command += " -Pattern \"" + f.pattern + "\""
	if f.quiet {
		f.command += " -quiet"
	}
	if f.caseSensitive {
		f.command += " -CaseSensitive"
	}

	output, err = interact.ShellPipe(f.command, s)

	if err != nil {
		return
	}

	artifact.SaveString(output, *s)

	return

}
