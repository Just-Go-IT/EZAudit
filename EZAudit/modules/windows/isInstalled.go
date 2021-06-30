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
	registry.Register("isInstalled", &isInstalled{}, false, registry.Windows)
}

type isInstalled struct {
	name    string
	command string
}

func (i isInstalled) New(p map[string]interface{}) (global.Module, error) {
	// Parse arguments
	ok := false
	// Check for optional parameter Mode
	for k := range p {
		switch k {
		case "name":
			i.name, ok = p["name"].(string)
			if !ok {
				return nil, errors.New("the key \"name\" is set for the module. The value must be a \"string\"")
			}
		default:
			return nil, errors.New("the key \"" + k + "\" is not supported in the module isInstalled")
		}

	}
	return &i, nil
}

func (i *isInstalled) Execute(s *global.Step) (output string, err error) {
	//Lock because of the mainResource

	global.IsInstalledLock.Lock()
	defer global.IsInstalledLock.Unlock()

	// build command
	i.command = "Get-ItemProperty HKLM:\\Software\\Wow6432Node\\Microsoft\\Windows\\CurrentVersion\\Uninstall\\* | Select-Object DisplayName, DisplayVersion, Publisher, InstallDate | Format-Table -AutoSize"

	//Check if isInstalled artifact already exists
	if _, errIsInstalledFile := os.Stat(global.IsInstalledPath); os.IsNotExist(errIsInstalledFile) {

		// path does not exist
		artifact.CreateFolderIfNotExist(global.ArtifactsDir)
		artifact.CreateFolderIfNotExist(global.MainResourcesDir)
		artifact.CreateFolderIfNotExist(global.IsInstalledDir)

		file, errC := os.Create(global.IsInstalledPath)
		if errC != nil {
			return "", errC
		}

		// interact command
		output, err = interact.ShellPipe(i.command, s)
		if err != nil {
			return output, err
		}

		_, err = file.WriteString(output)
		if err != nil {
			return
		}

		errClose := file.Close()
		if errClose != nil {
			log.Body("Whilst trying to close "+file.Name()+" this Error occurred:"+errClose.Error(), s.Path, log.Error)
		}
	}

	// artifact already exists
	var isInstalledFile []byte
	isInstalledFile, err = ioutil.ReadFile(global.IsInstalledPath)
	if err != nil {
		return
	}

	output = string(isInstalledFile)

	result := ""
	//Zeilenweise  das artifact ablegen wenn es enthalten ist
	if i.name != "" {
		lines := strings.Split(output, "\n")
		for _, line := range lines {
			if strings.Contains(strings.ToLower(line), i.name) {
				result += lines[0]
				result += line
			}
		}
	}

	artifact.SaveString(result, *s)

	return
}
