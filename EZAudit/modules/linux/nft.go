package linux

import (
	"Just-Go-IT/EZAudit/artifact"
	"Just-Go-IT/EZAudit/global"
	"Just-Go-IT/EZAudit/interact"
	"Just-Go-IT/EZAudit/registry"
)

func init() {
	registry.Register("nft", &nft{}, false, registry.Linux)
}

type nft struct {
	command string
}

func (n nft) New(p map[string]interface{}) (global.Module, error) {
	return &n, nil
}

func (n *nft) Execute(s *global.Step) (output string, err error) {
	n.command += "nft list ruleset"

	output, err = interact.ShellPipe(n.command, s)
	if err != nil {
		return
	}

	artifact.SaveString(output, *s)
	return
}
