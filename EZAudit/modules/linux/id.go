package linux

import (
	"Just-Go-IT/EZAudit/artifact"
	"Just-Go-IT/EZAudit/global"
	"Just-Go-IT/EZAudit/interact"
	"Just-Go-IT/EZAudit/registry"
)

func init() {
	registry.Register("id", &id{}, false, registry.Linux)
}

type id struct {
	command string
}

func (i id) New(p map[string]interface{}) (global.Module, error) {
	return &i, nil
}

func (i *id) Execute(s *global.Step) (output string, err error) {
	i.command += "id -u"

	output, err = interact.ShellPipe(i.command, s)
	if err != nil {
		return
	}

	artifact.SaveString(output, *s)

	return
}
