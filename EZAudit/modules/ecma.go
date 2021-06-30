package modules

import (
	"Just-Go-IT/EZAudit/global"
	"Just-Go-IT/EZAudit/log"
	"Just-Go-IT/EZAudit/registry"
	"errors"
	"github.com/dop251/goja"
	"io/ioutil"
	"strconv"
)

func init() {
	registry.Register("ecma", &ecma{}, false, registry.Linux|registry.Windows)
}

var artifactCounter int

type ecma struct {
	path     string
	script   string
	compiled *goja.Program
}

func (e *ecma) New(p map[string]interface{}) (global.Module, error) {
	ok := false

	// Checks for optional parameters and if the keys are supported
	for key := range p {
		switch key {
		case "path":
			e.path, ok = p["path"].(string)
			if !ok {
				return nil, errors.New("the parameter path must be a string")
			}
		case "script":
			e.script, ok = p["script"].(string)
			if !ok {
				return nil, errors.New("the parameter script must be a string")
			}
		default:
			if key != "path" && key != "script" {
				return nil, errors.New("the key " + key + " does not exists")
			}
		}
	}

	if e.path != "" && e.script != "" {
		return nil, errors.New("incorrect configuration, path and script is set. Only one can be chosen")
	}

	if e.path != "" {
		file, err := ioutil.ReadFile(e.path)

		if err != nil {
			return nil, errors.New("While opening the path " + e.path + " the following error occurred:\n" + err.Error())
		}

		e.compiled, err = goja.Compile("myProgram", string(file), false)
		if err != nil {
			return nil, errors.New("syntax of the ecma Script is wrong. Can't be compiled" + err.Error())
		}

	} else {
		var err error

		e.compiled, err = goja.Compile("myProgram", e.script, false)
		if err != nil {
			return nil, errors.New("syntax of the ecma Script is wrong. Can't be compiled" + err.Error())
		}
	}

	return e, nil
}

func (e *ecma) Execute(s *global.Step) (output string, err error) {
	vm := goja.New()
	vm.Set("execute", func(moduleName string, parameter map[string]interface{}, saveArtifact bool, usePipe bool) string {
		return executeModule(moduleName, parameter, saveArtifact, usePipe, *s)
	})
	vm.Set("result", func(result string) {
		output = result
	})
	_, err = vm.RunProgram(e.compiled)
	if err != nil {
		log.Body(err.Error(), s.Path, log.Error)
	}
	return
}

func executeModule(moduleName string, parameter map[string]interface{}, saveArtifact bool, usePipe bool, step global.Step) string {
	if moduleTemplate := registry.GetModule(moduleName); moduleTemplate != nil {
		//Create Step to pass it to the module
		stepEcma := global.Step{}
		stepEcma.ModuleName = moduleName
		stepEcma.DontSaveArtifact = saveArtifact
		stepEcma.Parameters = parameter
		stepEcma.UsePipe = usePipe
		stepEcma.Path = global.Path{
			CommandIndex: stepEcma.Path.CommandIndex,
			CommandName:  step.Path.CommandName,
			StepIndex:    stepEcma.Path.StepIndex,
			ExactPath:    stepEcma.Path.ExactPath + "_ScriptedStep_" + strconv.Itoa(artifactCounter),
		}

		artifactCounter++
		//Create module and execute it if it can be parsed
		if m, err := moduleTemplate.New(stepEcma.Parameters); err == nil {
			stepEcma.Module = m
			res, errE := stepEcma.Module.Execute(&stepEcma)
			if errE != nil {
				log.Body(errE.Error(), step.Path, log.Error)
			}
			return res
		} else {
			log.Body(err.Error(), step.Path, log.Error)
		}
	} else {
		log.Body("module "+moduleName+" not found in the registry", step.Path, log.Error)

	}
	return ""
}
