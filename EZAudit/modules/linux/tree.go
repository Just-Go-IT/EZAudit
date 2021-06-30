package linux

import (
	"Just-Go-IT/EZAudit/artifact"
	"Just-Go-IT/EZAudit/global"
	"Just-Go-IT/EZAudit/interact"
	"Just-Go-IT/EZAudit/registry"
)

func init() {
	registry.Register("tree", &tree{}, false, registry.Linux)
}

type tree struct {
	command string
}

func (t tree) New(p map[string]interface{}) (global.Module, error) {
	return &t, nil
}

func (t *tree) Execute(s *global.Step) (output string, err error) {
	t.command += "tree -a --dirsfirst -n"

	output, err = interact.ShellPipe(t.command, s)
	if err != nil {
		return
	}

	artifact.SaveString(output, *s)
	return
}
