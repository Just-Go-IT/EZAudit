// Package log implements all log function for the project.
package log

import (
	"Just-Go-IT/EZAudit/flags"
	"Just-Go-IT/EZAudit/global"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"sort"
	"strconv"
	"strings"
	"time"
)

const (
	Fatal       = 1
	Error       = 2
	Warning     = 3
	Information = 4
	Debug       = 5
)

var (
	toString = map[int]string{
		Fatal:       "FATAL:  ",
		Error:       "ERROR:  ",
		Warning:     "WARNING:",
		Information: "INFO:   ",
		Debug:       "DEBUG:  ",
	}
	commandIntToName map[int]string
	header           []message
	body             []map[string][]message // First array is CommandIndex, map-key is exact path and the values are the messages
	footer           []message
	verbosityLevel   = Information // Set the Default verbosityLevel to Information level
)

type message struct {
	message   string
	verbosity int
}

func setVerbosity(v int) {
	verbosityLevel = v
}

// SetVerbosityLevel Fatal = 1, Error = 2, Warning = 3, Information = 4, Debug = 5
func SetVerbosityLevel(v int) {
	switch v {
	case 1:
		setVerbosity(Fatal)
	case 2:
		setVerbosity(Error)
	case 3:
		setVerbosity(Warning)
	case 4:
		setVerbosity(Information)
	case 5:
		setVerbosity(Debug)
	default:
		headerSkipStackTrace("the verbosity Level "+strconv.Itoa(v)+" doesn't exist. Verbosity Level is set to "+strings.ReplaceAll(toString[verbosityLevel], ":", ""), Warning, 3)
	}
}

func CreateCommandLog(commandsLength int) {
	// Right length for Commands
	body = make([]map[string][]message, commandsLength)
	for i := range body {
		body[i] = make(map[string][]message)
	}
	commandIntToName = make(map[int]string)
}

// Write write the debug.log, but it writes only the messages which fulfill the verbosity. If the debug.log is empty the function will delete the debug.log.
func Write(path string) {

	// Create logfile
	logFile, err := os.Create(path)
	if err != nil {
		log.Fatal("Can't create logfile" + err.Error())
	}

	// Write header
	for _, msg := range header {
		if msg.verbosity <= verbosityLevel {
			logFile.WriteString(toString[msg.verbosity] + strings.ReplaceAll(msg.message, "\n", "\n\t\t\t") + "\n")
		}
	}

	// Write body (iterate commands)
	bodyIsEmpty := true
	for i, commands := range body {
		// If the message array is bodyIsEmpty, skip it
		if !commandHasMessageWithRightVerbosity(commands) {
			continue
		}
		bodyIsEmpty = false
		logFile.WriteString("########################################################################################################################\n")
		logFile.WriteString("Command " + strconv.Itoa(i) + "\t" + commandIntToName[i] + "\n")

		keys := make([]string, 0, len(commands))
		for k := range commands {
			keys = append(keys, k)
		}
		sort.Strings(keys)

		// Check verbosity level
		for _, k := range keys {
			if stepHasMessageWithRightVerbosity(commands[k]) {
				logFile.WriteString("\tStep " + k + " " + "\n")
				for _, msg := range commands[k] {
					if msg.verbosity <= verbosityLevel {
						logFile.WriteString("\t\t" + toString[msg.verbosity] + strings.ReplaceAll(msg.message, "\n", "\n\t\t\t\t\t") + "\n")
					}
				}
			}
		}

		// A separator line between body and footer for aesthetics
		if i == len(body)-1 && !bodyIsEmpty {
			logFile.WriteString("########################################################################################################################\n")
		}

	}

	// Write footer
	for _, msg := range footer {
		if msg.verbosity <= verbosityLevel {
			logFile.WriteString(toString[msg.verbosity] + strings.ReplaceAll(msg.message, "\n", "\n\t\t") + "\n")
		}
	}

	// Close file
	logFile.Close()

	// Delete if empty
	if stat, _ := os.Stat(path); stat.Size() == 0 {
		os.Remove(path)
	}
}

func commandHasMessageWithRightVerbosity(command map[string][]message) bool {
	for _, msgs := range command {
		for _, msg := range msgs {
			if msg.verbosity <= verbosityLevel {
				return true
			}
		}
	}
	return false
}

func stepHasMessageWithRightVerbosity(msgs []message) bool {
	for _, msg := range msgs {
		if msg.verbosity <= verbosityLevel {
			return true
		}
	}
	return false
}

// Header append a log message to the header
func Header(log string, verbosity int) {
	header = append(header, message{
		message:   getStackTrace(2) + log,
		verbosity: verbosity,
	})
}

func headerSkipStackTrace(log string, verbosity int, skips int) {
	header = append(header, message{
		message:   getStackTrace(skips) + log,
		verbosity: verbosity,
	})
}

// Body append a log message to the right command in the body based on the passed path
func Body(log string, path global.Path, verbosity int) {
	global.LoggerLock.Lock()
	bodySkipStackTrace(log, path, verbosity, 3)
	global.LoggerLock.Unlock()
}

func bodySkipStackTrace(log string, path global.Path, verbosity int, skips int) {
	global.CommandMapLock.Lock()
	if commandIntToName[path.CommandIndex] == "" {
		commandIntToName[path.CommandIndex] = path.CommandName
	}
	global.CommandMapLock.Unlock()
	m := body[path.CommandIndex]
	m[strconv.Itoa(path.StepIndex)+path.ExactPath] = append(m[strconv.Itoa(path.StepIndex)+path.ExactPath], message{message: getStackTrace(skips) + log, verbosity: verbosity})
}

// Footer append a log message to the header
func Footer(log string, verbosity int) {
	footer = append(footer, message{
		message:   getStackTrace(2) + "\t" + log,
		verbosity: verbosity,
	})
}

func getStackTrace(skip int) string {
	_, file, line, _ := runtime.Caller(skip)
	return time.Now().Format("2006-01-02 15:04:05") + " " + filepath.Base(file) + " " + strconv.Itoa(line) + ":\n"
}

// DryRunResult append a log message on the header based on the errorCounter
func DryRunResult(counter int) {
	if counter == 1 {
		headerSkipStackTrace(fmt.Sprintf("There is %v error in the config file."+
			" To see the required parameters and types for the specific module, please consult the manual", counter), Information, 3)
	} else if counter > 1 {
		headerSkipStackTrace(fmt.Sprintf("There are %v errors in the config file."+
			" To see the required parameters and types for the specific module, please consult the manual", counter), Information, 3)
	} else {
		headerSkipStackTrace("The config file is valid", Information, 3)
	}
}

// ModuleInfo append a log message to the body with a prettified version of the module
func ModuleInfo(module global.Module, path global.Path) {
	s := fmt.Sprintf("%#v", module)
	s = strings.Replace(s, "{", ":\n\t", 1)
	sub := strings.LastIndex(s, "}")
	s = s[:sub] + strings.Replace(s[sub:], "}", "", 1)
	s = strings.ReplaceAll(s, ", ", ",\n\t")
	Body("Module "+s, path, Debug)
}

// ModuleError append a log message to the body with a prettified version of the output and the error
func ModuleError(output string, err error, path global.Path) {
	s := "Module Error:\n\t" + strings.ReplaceAll(err.Error(), "\n", "\n\t")
	if output != "" {
		s += "\nModule Output:\n\t" + strings.ReplaceAll(output, "\n", "\n\t")
	}
	Body(s, path, Error)
}

// ModuleResult append a log message to the body with a prettified version of the auditStep
func ModuleResult(auditStep global.AuditStep, path global.Path) {
	s := fmt.Sprintf("%#v", auditStep)
	s = strings.Replace(s, "{", ":\n\t", 1)
	sub := strings.LastIndex(s, "}")
	s = s[:sub] + strings.Replace(s[sub:], "}", "", 1)
	s = strings.ReplaceAll(s, ", ", ",\n\t")
	Body(s, path, Debug)
}

// Summary append a prettified version of the summary to the footer
func Summary(summary global.Summary) {
	s := fmt.Sprintf("%#v", summary)
	s = strings.Replace(s, "{", ":\n\t\t", 1)
	sub := strings.LastIndex(s, "}")
	s = s[:sub] + strings.Replace(s[sub:], "}", "", 1)
	s = strings.ReplaceAll(s, ", ", ",\n\t\t")
	s = strings.ReplaceAll(s, ":", ": ")
	Footer(s, Information)
}

// DryRunError append log message to the body but only if the dryRun Flag is set
func DryRunError(log string, path global.Path) {
	if flags.IsSet("dryRun") {
		bodySkipStackTrace("Parsing: "+log, path, Error, 4)
	}
}

// Shell append log message to the body, if verbosity is Debug the log message is with more details
func Shell(command string, cmdStruct string, path *global.Path) {
	if verbosityLevel > Information {
		if path != nil {
			bodySkipStackTrace(fmt.Sprintf("Executed powershell command with all arguments:\n%v", cmdStruct), *path, Debug, 3)
		} else {
			headerSkipStackTrace(fmt.Sprintf("Executed powershell command with all arguments:\n%v", cmdStruct), Debug, 3)
		}
	} else {
		if path != nil {
			bodySkipStackTrace(fmt.Sprintf("Executed powershell command: %v", command), *path, Information, 3)
		} else {
			headerSkipStackTrace(fmt.Sprintf("Executed powershell command: %v", command), Information, 3)
		}
	}
}

// ActiveFlags append a log message to the header with all flags which are passed as arguments to the program
func ActiveFlags() {
	setFlags := flags.GetAllSetFlags()
	var flagNames []string
	for _, f := range setFlags {
		if f.GetValue() == "" {
			flagNames = append(flagNames, f.GetName())
		} else {
			flagNames = append(flagNames, f.GetName()+"="+f.GetValue())
		}
	}
	mode := strings.ReplaceAll(fmt.Sprintf("%v", flagNames), " ", ", ")
	mode = strings.ReplaceAll(mode, "[", "")
	mode = strings.ReplaceAll(mode, "]", "")

	if len(setFlags) == 0 {
		headerSkipStackTrace("Starting with no flags.", Information, 3)
	} else if len(setFlags) == 1 {
		headerSkipStackTrace("Starting with the flag: "+mode, Information, 3)
	} else {
		headerSkipStackTrace("Starting with flags: "+mode, Information, 3)
	}
}
