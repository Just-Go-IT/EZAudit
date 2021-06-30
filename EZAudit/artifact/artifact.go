// Package artifact implements functions for saving and manipulating strings.
package artifact

import (
	"Just-Go-IT/EZAudit/global"
	"Just-Go-IT/EZAudit/log"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
)

var (
	regex  []regexp.Regexp
	maxLen int
)

// SetUp compiles the passed regexExpressions and saves them for later use. If maxOuputLenght is defined in the config-file then it overrides the maxOutputLength (which is by default = 350).
// Set maxOutputLength to -1, if you want to set it to unlimited.
func SetUp(regexExpression []string, maxOutputLength int) {

	for i := range regexExpression {
		tmp, err := regexp.Compile(regexExpression[i])
		if err != nil {
			log.Header("Regex "+strconv.Itoa(i)+": "+regexExpression[i]+" couldn't be compiled:\n"+err.Error(), log.Error)
		} else {
			regex = append(regex, *tmp)
		}
	}
	if maxOutputLength == 0 {
		maxLen = 350
	} else {
		maxLen = maxOutputLength
	}
	log.Header("maxOutputLength is set to "+strconv.Itoa(maxLen), log.Information)
}

// CreateFolderIfNotExist creates an empty folder in the specified path, if the folder doesn't exist yet.
func CreateFolderIfNotExist(path string) {
	if _, errStat := os.Stat(path); os.IsNotExist(errStat) {
		err := os.Mkdir(path, global.FilePermissionValue)
		if err != nil {
			log.Header("The folder "+path+" couldn't be created. This error occurred:\n"+err.Error(), log.Error)
		}
	}
}

// SaveString saves string as a .txt File in the Artifact folder. The created fileName is based on getFilePath.
// The saved string depends on the censoring configuration and DontSaveArtifact parameter in the config-file.
func SaveString(artifact string, step global.Step) {

	if step.DontSaveArtifact || artifact == "" {
		return
	}

	//Creates the artifact folder
	CreateFolderIfNotExist(global.ArtifactsDir)
	CreateFolderIfNotExist(global.AuditDir)
	CreateFolderIfNotExist(getDirPath(step.Path))

	// Creates the resultFile
	if _, err := os.Stat(getFilePath(step.Path)); err == nil {
		//if file already exists
		step.Path.ExactPath += "_1"
		SaveString(artifact, step)
	} else if os.IsNotExist(err) {
		file, errC := os.Create(getFilePath(step.Path))
		if errC != nil {
			log.Header("The file  "+getFilePath(step.Path)+" couldn't be created. This error occurred:\n"+errC.Error(), log.Error)
		} else {
			// If file opens successfully, it closes it at the end
			defer func(file *os.File) {
				errClose := file.Close()
				if errClose != nil {
					log.Header("Whilst trying to close "+file.Name()+" this Error occurred:\n"+errClose.Error(), log.Error)
				}
			}(file)
		}

		// Sanitize the artifact
		artifact = CensorLocal(artifact, step.Regex)
		artifact = CensorGlobal(artifact)

		// Write the artifact into a file
		_, err = file.WriteString(artifact)
		if err != nil {
			log.Header("Couldn't write the artifact to this file "+getFilePath(step.Path)+". This error occurred:\n"+err.Error(), log.Error)
		}
		log.Body("Artifact saved to "+getFilePath(step.Path), step.Path, log.Information)
	} else if err != nil {
		log.Header("Couldn't create a artifact for this file "+getFilePath(step.Path)+". This error occurred:\n"+err.Error(), log.Error)
	}
}

// SaveFile save File to the Artifact folder. The created fileName is based on getFilePath.
// The saved string depends on the censoring configuration and DontSaveArtifact parameter in the config-file.
func SaveFile(path string, step global.Step) {
	if step.DontSaveArtifact || path == "" {
		return
	}

	// Create the artifact folder
	CreateFolderIfNotExist(global.ArtifactsDir)
	CreateFolderIfNotExist(global.AuditDir)
	CreateFolderIfNotExist(getDirPath(step.Path))

	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		log.Body("The artifact path "+path+" doesn't exists.\n"+err.Error(), step.Path, log.Error)
	} else {
		// Open the artifact source
		tmp, errF := ioutil.ReadFile(path)
		if errF != nil {
			log.Body("Whilst tyring to open the file "+path+" this failure occurred:\n"+errF.Error(), step.Path, log.Error)
		}

		source := string(tmp)

		// Sanitize the artifact
		source = CensorLocal(source, step.Regex)
		source = CensorGlobal(source)

		// Create a file for the artifact
		if _, errR := os.Stat(getFilePath(step.Path)); errR == nil {
			//file already exists
			step.Path.ExactPath += "_1"
			SaveFile(path, step)
		} else if os.IsNotExist(errR) {
			destination, errC := os.Create(getFilePath(step.Path))
			if errC != nil {
				log.Body("The file "+getFilePath(step.Path)+" couldn't be created. This error occurred:\n"+errC.Error(), step.Path, log.Error)
			}
			//Close File and handle error
			defer func(destination *os.File) {
				errClose := destination.Close()
				if errClose != nil {
					log.Body("Whilst trying to close "+getFilePath(step.Path)+" this Error occurred:\n"+errClose.Error(), step.Path, log.Error)
				}
			}(destination)

			// Writes it down
			_, errW := destination.WriteString(source)
			if errW != nil {
				log.Header("Couldn't write the artifact into the file "+path+". This error occurred:\n"+errW.Error(), log.Error)
			}
			log.Body("Artifact saved to "+getFilePath(step.Path), step.Path, log.Information)
		} else {
			log.Header("Couldn't create a artifact for this file "+getFilePath(step.Path)+". This error occurred:\n"+err.Error(), log.Error)
		}
	}
}

// getDirPath returns directory path based on global.Path
func getDirPath(p global.Path) string {
	return filepath.Join(global.AuditDir, "_"+fmt.Sprintf("%03d", p.CommandIndex)+"_"+sanitizeDirName(p.CommandName))
}

// getFilePath returns file path based on global.Path
func getFilePath(p global.Path) string {
	return filepath.Join(getDirPath(p), strconv.Itoa(p.StepIndex)+p.ExactPath+".txt")
}

// CheckLength checks if the length is too long, if it is it, it saves an artifact with the input and returns the path to it
func CheckLength(input string, step global.Step) string {
	if maxLen == -1 {
		return input
	}
	if maxLen < len(input) {
		if _, err := os.Stat(getFilePath(step.Path)); os.IsExist(err) {
			return getFilePath(step.Path)
		} else {
			SaveString(input, step)
			return getFilePath(step.Path)
		}
	}
	return input
}

// CensorGlobal takes the string and replaces the text, based on the global censor regexp to "--censored--".
// Returns the sanitized string
func CensorGlobal(input string) (censored string) {
	for i := range regex {
		input = regex[i].ReplaceAllString(input, "--censored--")
	}
	return input
}

// CensorLocal takes string and specific regexp. Replaces text based on the specific regexp to "--censored--".
// Returns the sanitized string
func CensorLocal(input string, reg []regexp.Regexp) string {
	for i := range reg {
		input = reg[i].ReplaceAllString(input, "--censored--")
	}
	return input
}

// sanitizeDirName sanitizes unsupported File-System characters
func sanitizeDirName(s string) string {
	// Start with lowercase string
	fileName := strings.ToLower(s)
	fileName = path.Clean(path.Base(fileName))

	// Remove all unsupported characters for directory names, replacing them with a _
	return sanitizeString(fileName, regexp.MustCompile(`[^[:alnum:]-.]`))
}

// sanitizeString sanitizes unsupported File-System characters
func sanitizeString(s string, r *regexp.Regexp) string {

	// Remove any trailing whitespace
	s = strings.Trim(s, " ")

	// Replace certain dangerous characters with a underscore
	s = regexp.MustCompile(`[ &_=+:]`).ReplaceAllString(s, "_")

	// Remove all other unrecognised characters based on the passed regex
	s = r.ReplaceAllString(s, "")

	// Clean up double underscores
	return regexp.MustCompile(`[\-]+`).ReplaceAllString(s, "_")
}
