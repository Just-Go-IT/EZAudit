package linux

import (
	"Just-Go-IT/EZAudit/artifact"
	"Just-Go-IT/EZAudit/global"
	"Just-Go-IT/EZAudit/interact"
	"Just-Go-IT/EZAudit/registry"
	"errors"
)

func init() {
	registry.Register("whereis", &whereIs{}, false, registry.Linux)
}

type whereIs struct {
	target  string
	command string
}

func (w whereIs) New(p map[string]interface{}) (global.Module, error) {
	ok := false

	w.target, ok = p["target"].(string)
	if !ok {
		return nil, errors.New("the key \"target\" must be set and the value must be a \"string\"")
	}
	if w.target == "" {
		return nil, errors.New("there is no action option in the module whereis")
	}
	return &w, nil
}

func (w *whereIs) Execute(s *global.Step) (output string, err error) {

	w.command += "whereis -b  " + w.target

	output, err = interact.ShellPipe(w.command, s)

	if err != nil {
		return
	}

	artifact.SaveString(output, *s)
	return
}
