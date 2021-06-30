package linux

import (
	"Just-Go-IT/EZAudit/artifact"
	"Just-Go-IT/EZAudit/global"
	"Just-Go-IT/EZAudit/interact"
	"Just-Go-IT/EZAudit/registry"
	"errors"
)

func init() {
	registry.Register("rpm", &rpm{}, false, registry.Linux)
}

type rpm struct {
	target   string
	query    bool
	queryAll bool
	verify   bool
	command  string
}

func (r rpm) New(p map[string]interface{}) (global.Module, error) {

	ok := false

	// Checks for optional parameters and if the keys are supported
	for key := range p {
		switch key {
		case "target":
			r.target, ok = p["target"].(string)
			if !ok {
				return nil, errors.New("the key \"target\" is set for the module. The value must be a \"string\"")
			}
		case "query":
			r.query, ok = p["query"].(bool)
			if !ok {
				return nil, errors.New("the key \"query\" is set for the module. The value must be a \"bool\"")
			}
		case "queryAll":
			r.queryAll, ok = p["queryAll"].(bool)
			if !ok {
				return nil, errors.New("the key \"queryAll\" is set for the module. The value must be a \"bool\"")
			}
		case "verify":
			r.verify, ok = p["verify"].(bool)
			if !ok {
				return nil, errors.New("the key \"verify\" is set for the module. The value must be a \"bool\"")
			}
		default:
			if key != "target" && key != "query" && key != "queryAll" && key != "verify" {
				return nil, errors.New("there is no key called: \"" + key + "\" in the module rpm")
			}
		}
	}
	if (!r.verify && !r.query && !r.queryAll && r.target != "") || (r.query && r.verify && r.queryAll) {
		return nil, errors.New("there is no action option in the module rpm")
	}

	return &r, nil
}

func (r *rpm) Execute(s *global.Step) (output string, err error) {
	r.command += "rpm"

	if r.query && r.target != "" {
		r.command += " -q " + r.target
	}
	if r.verify && r.target != "" {
		r.command += " -Va " + r.target
	}
	if r.queryAll && r.target != "" {
		r.command += " -pq " + r.target
	}

	output, err = interact.ShellPipe(r.command, s)

	if err != nil {
		return
	}

	artifact.SaveString(output, *s)
	return
}
