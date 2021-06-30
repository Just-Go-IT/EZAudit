package linux

import (
	"Just-Go-IT/EZAudit/artifact"
	"Just-Go-IT/EZAudit/global"
	"Just-Go-IT/EZAudit/interact"
	"Just-Go-IT/EZAudit/registry"
)

func init() {
	registry.Register("useradd", &useradd{}, false, registry.Linux)
}

type useradd struct {
	command string
}

func (u useradd) New(p map[string]interface{}) (global.Module, error) {

	return &u, nil
}

func (u *useradd) Execute(s *global.Step) (output string, err error) {
	u.command += "useradd-D "

	output, err = interact.ShellPipe(u.command, s)

	if err != nil {
		return
	}

	artifact.SaveString(output, *s)
	return
}
