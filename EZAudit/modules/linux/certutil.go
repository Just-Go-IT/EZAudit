package linux

import (
	"Just-Go-IT/EZAudit/artifact"
	"Just-Go-IT/EZAudit/global"
	"Just-Go-IT/EZAudit/interact"
	"Just-Go-IT/EZAudit/registry"
)

func init() {
	registry.Register("certutil", &certutil{}, false, registry.Linux)
}

type certutil struct {
	command string
}

func (c certutil) New(p map[string]interface{}) (global.Module, error) {
	return &c, nil
}

func (c *certutil) Execute(s *global.Step) (output string, err error) {
	c.command = "certutil-L"

	output, err = interact.ShellPipe(c.command, s)
	if err != nil {
		return
	}

	artifact.SaveString(output, *s)

	return
}
