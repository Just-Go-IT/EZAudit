package linux

import (
	"Just-Go-IT/EZAudit/artifact"
	"Just-Go-IT/EZAudit/global"
	"Just-Go-IT/EZAudit/interact"
	"Just-Go-IT/EZAudit/registry"
)

func init() {
	registry.Register("appArmour", &appArmor{}, true, registry.Linux)
}

type appArmor struct {
	command string
}

func (app appArmor) New(p map[string]interface{}) (global.Module, error) {
	return &app, nil
}

func (app *appArmor) Execute(s *global.Step) (output string, err error) {
	app.command += "apparmor_status "

	output, err = interact.ShellPipe(app.command, s)
	if err != nil {
		return
	}

	artifact.SaveString(output, *s)

	return
}
