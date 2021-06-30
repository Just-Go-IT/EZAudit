// Package interact provides functions to interact with other programs or the shell.
package interact

import (
	"Just-Go-IT/EZAudit/global"
	"Just-Go-IT/EZAudit/log"
	"fmt"
	"os/exec"
	"strings"
)

// Thread-safe, each Body has his own pipe
var pipeStorage map[int][]string

func init() {
	pipeStorage = make(map[int][]string)
}

// Program  executes the passed program with the arguments and returns stdOut and stdErr as a string and an error. If no error happened, the error is nil
func Program(s *global.Step, programName string, arguments ...string) (combinedOutput string, err error) {
	cmd := exec.Command(programName, arguments...)
	log.Body(fmt.Sprintf("Executed program: %v with the arguments %v", programName, arguments), s.Path, log.Information)
	cmdOut, err := cmd.CombinedOutput()
	combinedOutput = strings.TrimSpace(string(cmdOut))
	return
}

// Shell executes the passed command in the shell and returns stdOut and stdErr as a string and an error. If no error happened, the error is nil
func Shell(command string) (combinedOutput string, err error) {
	cmd := exec.Command("powershell.exe", "-NoLogo", "-NoProfile", "-NonInteractive", command)
	log.Shell(command, cmd.String(), nil)
	cmdOut, err := cmd.CombinedOutput()
	combinedOutput = strings.TrimSpace(string(cmdOut))
	return
}

// ShellPipe executes the passed command in the shell with a pipe if usePipe is true and returns stdOut and stdErr as a string and an error. If no error happened, the error is nil. Other to Shell it saves the command if the next Module needs a pipe and delete the pipe if the next Module don't need a pipe.
func ShellPipe(cmdIn string, s *global.Step) (combinedOutput string, err error) {
	// next step needs a pipe
	nextStepNeedsPipe := s.OnSuccess != nil && s.OnSuccess.UsePipe || s.OnFailure != nil && s.OnFailure.UsePipe ||
		s.GetNext != nil && s.GetNext.UsePipe
	if s.UsePipe {

		// We need to build them new every time because a exec.Cmd can only be executed once
		var cmd string

		for _, commandString := range pipeStorage[s.Path.CommandIndex] {
			cmd += commandString + " | "
		}
		cmd += cmdIn

		command := exec.Command("powershell.exe", "-NoProfile", "-NoLogo", "-NonInteractive", cmd)
		log.Shell(cmd, command.String(), &s.Path)
		var tmp []byte

		tmp, err = command.CombinedOutput()
		combinedOutput = strings.TrimSpace(string(tmp))

	}

	if nextStepNeedsPipe {
		// Prepare pipe for the next step
		global.PipeStorageLock.Lock()
		pipeStorage[s.Path.CommandIndex] = append(pipeStorage[s.Path.CommandIndex], cmdIn)
		global.PipeStorageLock.Unlock()
	}

	if s.UsePipe && !nextStepNeedsPipe {
		// Clean the command pipe
		global.PipeStorageLock.Lock()
		delete(pipeStorage, s.Path.CommandIndex)
		global.PipeStorageLock.Unlock()

	}

	if !s.UsePipe {
		var tmp []byte
		command := exec.Command("powershell.exe", "-NoProfile", "-NoLogo", "-NonInteractive", cmdIn)
		log.Shell(cmdIn, command.String(), &s.Path)
		tmp, err = command.CombinedOutput()
		combinedOutput = strings.TrimSpace(string(tmp))
	}

	return
}

// IsElevated The manifest forces the program to be executed as admin
func IsElevated() bool {
	return true
}
