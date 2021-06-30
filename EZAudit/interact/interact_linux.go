// Package interact provides functions to interact with other programs or the shell.
package interact

import (
	"Just-Go-IT/EZAudit/global"
	"Just-Go-IT/EZAudit/log"
	"fmt"
	"os/exec"
	"strconv"
	"strings"
)

// Thread-safe, each Body has its own pipe
var pipeStorage map[int][]string

func init() {
	pipeStorage = make(map[int][]string)
}

// Program executes the passed program with the arguments and returns stdOut and stdErr as a string and an error. If no error happened, the error is nil.
func Program(s *global.Step, programName string, arguments ...string) (combinedOutput string, err error) {
	cmd := exec.Command(programName, arguments...)
	log.Body(fmt.Sprintf("Executed program: %v with the arguments %v", programName, arguments), s.Path, log.Information)
	cmdOut, err := cmd.CombinedOutput()
	combinedOutput = strings.TrimSpace(string(cmdOut))
	return
}

// Shell executes the passed command in the shell and returns stdOut and stdErr as a string and an error. If no error happened, the error is nil
func Shell(command string) (combinedOutput string, err error) {
	cmd := exec.Command("bash", "-c", command)
	log.Shell(command, cmd.String(), nil)
	cmdOut, err := cmd.CombinedOutput()
	combinedOutput = strings.TrimSpace(string(cmdOut))
	return
}

// ShellPipe executes the passed command in the shell with a pipe if usePipe is true and returns stdOut and stdErr as a string and an error. If no error happened, the error is nil. Other to Shell it saves the command if the next Module needs a pipe and delete the pipe if the next Module don't need a pipe.
func ShellPipe(cmd string, s *global.Step) (combinedOutput string, err error) {
	// Next step needs a pipe
	nextStepNeedsPipe := s.OnSuccess != nil && s.OnSuccess.UsePipe || s.OnFailure != nil && s.OnFailure.UsePipe ||
		s.GetNext != nil && s.GetNext.UsePipe

	command := exec.Command("bash", "-c", cmd)
	if s.UsePipe {

		// We need to build them new every time because the exec.Cmd can only be executed once
		var commands []*exec.Cmd
		logMsg := ""
		logCmd := ""
		for _, commandString := range pipeStorage[s.Path.CommandIndex] {
			c := exec.Command("bash", "-c", commandString)
			logMsg += commandString + " | "
			logCmd += c.String() + " | "
			commands = append(commands, c)
		}

		// Connect all commands through a pipe
		for i := range commands {
			pipe, _ := commands[i].StdoutPipe()

			// Last command needs to be piped into the actual command
			if i == len(pipeStorage[s.Path.CommandIndex])-1 {

				command.Stdin = pipe
				commands[i].Start()
				var tmp []byte
				logCmd += command.String()
				logMsg += cmd
				log.Shell(logMsg, logCmd, &s.Path)
				tmp, err = command.CombinedOutput()
				combinedOutput = strings.TrimSpace(string(tmp))
				break
			}

			commands[i+1].Stdin = pipe
			commands[i].Start()
		}
	}

	if nextStepNeedsPipe {
		// Prepare the pipe for the next step
		global.PipeStorageLock.Lock()
		pipeStorage[s.Path.CommandIndex] = append(pipeStorage[s.Path.CommandIndex], cmd)
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
		log.Shell(cmd, command.String(), &s.Path)
		tmp, err = command.CombinedOutput()
		combinedOutput = strings.TrimSpace(string(tmp))
	}
	return
}

// IsElevated checks if the program is elevated
func IsElevated() bool {
	cmd := exec.Command("id", "-u")
	output, err := cmd.Output()

	if err != nil {
		fmt.Println(err.Error())
	}

	// 0 = root, 501 = non-root user
	i, err := strconv.Atoi(string(output[:len(output)-1]))

	if err != nil {
		fmt.Println(err.Error())
	}

	return i == 0
}
