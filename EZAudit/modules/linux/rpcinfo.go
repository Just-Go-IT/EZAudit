package linux

import (
	"Just-Go-IT/EZAudit/artifact"
	"Just-Go-IT/EZAudit/global"
	"Just-Go-IT/EZAudit/interact"
	"Just-Go-IT/EZAudit/registry"
)

func init() {
	registry.Register("rpcinfo", &rpcinfo{}, false, registry.Linux)
}

type rpcinfo struct {
	command string
}

func (r rpcinfo) New(p map[string]interface{}) (global.Module, error) {
	return &r, nil
}

func (r *rpcinfo) Execute(s *global.Step) (output string, err error) {
	r.command += "rpcinfo -p"

	output, err = interact.ShellPipe(r.command, s)

	if err != nil {
		return
	}

	artifact.SaveString(output, *s)
	return
}
