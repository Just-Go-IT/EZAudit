// Package windows includes all modules which are usable for windows
package windows

import (
	"Just-Go-IT/EZAudit/artifact"
	"Just-Go-IT/EZAudit/global"
	"Just-Go-IT/EZAudit/interact"
	"Just-Go-IT/EZAudit/registry"
	"errors"
	"io/ioutil"
	"os"
	"strings"
)

func init() {
	registry.Register("auditpol", &auditpol{}, false, registry.Windows)
}

type auditpol struct {
	machineName      string
	policyTarget     string
	subCategory      string
	subCategoryGUID  string
	inclusionSetting string
	exclusionSetting string
}

func (a auditpol) New(p map[string]interface{}) (global.Module, error) {
	matchInterface, ok := p["match"].(map[string]interface{})
	if !ok {
		return nil, errors.New("the key \"match\" must be set and the value must be an \"object with key, value pairs\"")
	}
	for key, value := range matchInterface {
		valueOk := false
		switch key {
		case "machineName":
			a.machineName, valueOk = value.(string)
			if !valueOk {
				return nil, errors.New("the key \"machineName\" is set but the value is not a \"string\"")
			}
		case "policyTarget":
			a.policyTarget, valueOk = value.(string)
			if !valueOk {
				return nil, errors.New("the key \"policyTarget\" is set but the value is not a \"string\"")
			}
		case "subCategory":
			a.subCategory, valueOk = value.(string)
			if !valueOk {
				return nil, errors.New("the key \"subCategory\" is set but the value is not a \"string\"")
			}
		case "subCategoryGUID":
			a.subCategoryGUID, valueOk = value.(string)
			if !valueOk {
				return nil, errors.New("the key \"subCategoryGUID\" is set but the value is not a \"string\"")
			}
		case "inclusionSetting":
			a.inclusionSetting, valueOk = value.(string)
			if !valueOk {
				return nil, errors.New("the key \"inclusionSetting\" is set but the value is not a \"string\"")
			}
		case "exclusionSetting":
			a.exclusionSetting, valueOk = value.(string)
			if !valueOk {
				return nil, errors.New("the key \"exclusionSetting\" is set but the value is not a \"string\"")
			}
		default:
			return nil, errors.New("the key \"" + key + "\" does not exists in the module auditpol")
		}
	}
	return &a, nil
}

func (a *auditpol) Execute(s *global.Step) (output string, err error) {
	//Lock because of the mainResource
	global.AuditpolLock.Lock()
	defer global.AuditpolLock.Unlock()

	//Check if auditpol artifact already exists
	if _, errAuditpolFile := os.Stat(global.AuditpolPath); os.IsNotExist(errAuditpolFile) {
		// path does not exist
		artifact.CreateFolderIfNotExist(global.ArtifactsDir)
		artifact.CreateFolderIfNotExist(global.MainResourcesDir)
		artifact.CreateFolderIfNotExist(global.AuditpolDir)
		r, errA := interact.Program(s, "auditpol.exe", "/backup", "/file:"+strings.Replace(global.AuditpolPath, "/", "\\", -1)+"")
		if errA != nil {
			return r, errA
		}
		// path exists
		var file []byte
		file, err = ioutil.ReadFile(global.AuditpolPath)
		if err != nil {
			return
		}
		output = artifact.CensorGlobal(string(file))

		err = ioutil.WriteFile(global.AuditpolPath, []byte(output), global.FilePermissionValue)
		if err != nil {
			return
		}
	}

	// path exists
	var file []byte
	file, err = ioutil.ReadFile(global.AuditpolPath)
	if err != nil {
		return
	}

	output = string(file)

	//lines are the single objects
	lines := strings.Split(output, "\n")
	var matchedLines []string

	for _, line := range lines {
		//single entries in an object
		match := false

		values := strings.Split(line, ",")
		for j, value := range values {
			switch j {
			case 0:
				if a.machineName != "" {
					if a.machineName == value {
						match = true
					} else {
						match = false
						goto notPassed
					}
				}
			case 1:
				if a.policyTarget != "" {
					if a.policyTarget == value {
						match = true
					} else {
						match = false
						goto notPassed
					}
				}
			case 2:
				if a.subCategory != "" {
					if a.subCategory == value {
						match = true
					} else {

						match = false
						goto notPassed
					}
				}
			case 3:
				if a.subCategoryGUID != "" {
					if a.subCategoryGUID == value {
						match = true
					} else {
						match = false
						goto notPassed
					}
				}
			case 4:
				if a.inclusionSetting != "" {
					if a.inclusionSetting == value {
						match = true
					} else {
						match = false
						goto notPassed
					}
				}
			case 5:
				if a.exclusionSetting != "" {
					if a.exclusionSetting == value {
						match = true
					} else {
						match = false
						goto notPassed
					}
				}
			}
		}
	notPassed:
		if match {
			matchedLines = append(matchedLines, line)
		}

	}

	//How many matches do we have?
	if len(matchedLines) > 1 { //there are more than one match
		for _, matchedLine := range matchedLines {
			output = strings.ReplaceAll(strings.Split(matchedLine, ",")[6], "\r", "")
			if s.ExpectedValue != output { //if there is one result, that doesnt match the expectedValue than return it.
				//break and return the wrong output after saving the artifact
				break
			}
		}
	} else if len(matchedLines) == 1 { //there is only one match
		output = strings.ReplaceAll(strings.Split(matchedLines[0], ",")[6], "\r", "")
	}

	artifact.SaveString(strings.Join(matchedLines, "\n"), *s)

	return output, nil
}
