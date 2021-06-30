package windows

import (
	"Just-Go-IT/EZAudit/artifact"
	"Just-Go-IT/EZAudit/global"
	"Just-Go-IT/EZAudit/interact"
	"Just-Go-IT/EZAudit/log"
	"Just-Go-IT/EZAudit/registry"
	"errors"
	"io/ioutil"
	"os"
	"strings"
)

func init() {
	registry.Register("windowsServices", &windowsServices{}, false, registry.Windows)
}

type windowsServices struct {
	name    string
	status  string
	command string
}

func (s windowsServices) New(p map[string]interface{}) (global.Module, error) {
	// Parse arguments
	ok := false
	// Check for optional parameter Mode
	for k := range p {
		switch k {
		case "name":
			s.name, ok = p["name"].(string)
			if !ok {
				return nil, errors.New("the key \"name\" is set for the module. The value must be a \"string\"")
			}
		case "status":
			s.status, ok = p["status"].(string)
			if !ok {
				return nil, errors.New("the key \"status\" is set for the module. The value must be a \"string\"")
			}
		default:
			return nil, errors.New("the key \"" + k + "\" is not supported in the module isInstalled")
		}

	}
	return &s, nil
}

func (s *windowsServices) Execute(step *global.Step) (output string, err error) {
	//Lock because of the mainResource
	global.WindowsServicesLock.Lock()
	defer global.WindowsServicesLock.Unlock()

	// build command
	s.command += "Get-Service | Format-Table -AutoSize"

	//Check if isInstalled artifact already exists
	if _, errWindowsServicePath := os.Stat(global.WindowsServicesPath); os.IsNotExist(errWindowsServicePath) {
		// artifact does not exist
		artifact.CreateFolderIfNotExist(global.ArtifactsDir)
		artifact.CreateFolderIfNotExist(global.MainResourcesDir)
		artifact.CreateFolderIfNotExist(global.IsInstalledDir)

		file, errF := os.Create(global.WindowsServicesPath)
		if errF != nil {
			return "", errF
		}

		// interact command
		output, err = interact.ShellPipe(s.command, step)
		if err != nil {
			return output, err
		}
		_, err = file.WriteString(output)
		if err != nil {
			return output, err
		}

		errClose := file.Close()
		if errClose != nil {
			log.Body("Whilst trying to close "+file.Name()+" this Error occurred:\n"+errClose.Error(), step.Path, log.Error)
		}

	}
	// path exists
	file, errF := ioutil.ReadFile(global.WindowsServicesPath)
	if err != nil {
		return "", errF
	}
	output = string(file)

	result := ""
	lines := strings.Split(output, "\n")

	if s.name != "" {
		for _, line := range lines {
			if strings.Contains(strings.ToLower(line), s.name) {
				if s.status != "" {
					if strings.Contains(strings.ToLower(line), strings.ToLower(s.status)) {
						result += line
					}
				} else {
					result += line
				}
			}
		}
	}

	if s.status != "" {
		for _, line := range lines {
			if strings.Contains(strings.ToLower(line), strings.ToLower(s.status)) {
				if s.name != "" {
					if strings.Contains(strings.ToLower(line), strings.ToLower(s.name)) {
						if !strings.Contains(result, line) {
							result += line
						}
					}
				} else {
					result += line
				}
			}
		}
	}

	artifact.SaveString(result, *step)
	return
}
