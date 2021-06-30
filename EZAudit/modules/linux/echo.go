package linux

import (
	"Just-Go-IT/EZAudit/artifact"
	"Just-Go-IT/EZAudit/global"
	"Just-Go-IT/EZAudit/interact"
	"Just-Go-IT/EZAudit/registry"
	"errors"
)

func init() {
	registry.Register("echo", &echo{}, false, registry.Linux)
}

type echo struct {
	command string
	target  string
}

func (e echo) New(p map[string]interface{}) (global.Module, error) {
	ok := false

	e.target, ok = p["target"].(string)
	if !ok {
		return nil, errors.New("the key \"target\" must be set and the value must be a \"string\"")
	}

	if e.target == "" {
		return nil, errors.New("there is no action option in the module echo")
	}

	return &e, nil
}

func (e *echo) Execute(currentStep *global.Step) (output string, err error) {
	e.command += "echo"

	if e.target != "" {
		e.command += " " + e.target
	}

	output, err = interact.ShellPipe(e.command, currentStep)
	if err != nil {
		return
	}

	artifact.SaveString(output, *currentStep)
	return
}
