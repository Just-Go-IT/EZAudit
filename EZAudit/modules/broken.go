package modules

import (
	"Just-Go-IT/EZAudit/global"
	"Just-Go-IT/EZAudit/registry"
	"errors"
)

func init() {
	registry.Register("brokenModule", BrokenModule{}, false, registry.SupportAll)
}

type BrokenModule struct {
	Reason string
}

func (b BrokenModule) New(parameters map[string]interface{}) (global.Module, error) {
	return &b, nil
}

func (b BrokenModule) Execute(currentStep *global.Step) (output string, err error) {
	return "", errors.New(b.Reason)
}
