// Package executer implements functions for converting commands to audits and evaluate the result.
package executer

import (
	"Just-Go-IT/EZAudit/artifact"
	"Just-Go-IT/EZAudit/global"
	"Just-Go-IT/EZAudit/log"
	"fmt"
	"github.com/schollz/progressbar/v3"
	"strconv"
	"strings"
	"sync"
)

const (
	success = "success"
	fail    = "fail"
	error   = "error"
)

// ExecuteAllCommands processes all given Commands to Audits, so it executes the modules in the specific steps and evaluates the results
func ExecuteAllCommands(commands *[]global.Command) (audits []global.Audit) {

	fmt.Println("Executing all Audits...")

	commandAmount := len(*commands)

	audits = make([]global.Audit, commandAmount)
	progressBar := progressbar.New(commandAmount)
	var waitGroup sync.WaitGroup
	waitGroup.Add(commandAmount)

	// iterates the Commands
	for i := range *commands {
		audits[i].ID = i
		audits[i].Name = (*commands)[i].Name
		go executeCommand(&(*commands)[i], &audits[i], &waitGroup, progressBar)
	}

	// Waits until all go routines are finished
	waitGroup.Wait()

	return audits
}

// executeCommand processes a given command to an audit, so it executes the modules in the command and evaluates the result
func executeCommand(command *global.Command, audit *global.Audit, wg *sync.WaitGroup, bar *progressbar.ProgressBar) {
	defer wg.Done()
	defer bar.Add(1)

	audit.AuditSteps = make([]global.AuditStep, len(command.Steps))

	// Iterating Steps inside the Body
	for i := range command.Steps {
		// Execute all Modules in one step and get the auditSteps back
		audit.AuditSteps[i] = *executeStep(&command.Steps[i])
	}

	// Evaluate the audit status
	for i := range audit.AuditSteps {
		if audit.AuditSteps[i].Status == error {
			audit.Status = error
			break
		}
		if audit.AuditSteps[i].Status == fail && !command.Steps[i].AllowFailure {
			audit.Status = fail
			break
		}
		audit.Status = success
	}
}

// executeStep processes a given step to a audit, so it executes all modules in the step and evaluate the results
func executeStep(step *global.Step) (auditStep *global.AuditStep) {
	auditStep = &global.AuditStep{}
	auditStep.ID = step.Path.StepIndex
	auditStep.Expected = step.ExpectedValue
	auditStep.Comparison = step.Comparison

	log.ModuleInfo(step.Module, step.Path)

	// Execute the Module and get the Results
	output, err := step.Module.Execute(step)

	// If the module runs into an error
	if err != nil {
		output = artifact.CensorGlobal(output)
		output = artifact.CensorLocal(output, step.Regex)
		output = artifact.CheckLength(output, *step)
		auditStep.Status = error
		log.ModuleError(output, err, step.Path)
		return auditStep
	}

	// Check the Results and return if the step was successful and its description
	auditStep.Status, auditStep.Description = evaluateResult(output, step.Comparison, step.ExpectedValue)
	output = artifact.CensorGlobal(output)
	output = artifact.CensorLocal(output, step.Regex)
	auditStep.RealValue = artifact.CheckLength(output, *step)

	log.ModuleResult(*auditStep, step.Path)
	if auditStep.Description != "" && auditStep.Status == error {
		log.Body(auditStep.Description, step.Path, log.Error)
	}

	// Recursion
	if auditStep.Status == success {
		if step.OnSuccess != nil {
			auditStep.OnSuccess = executeStep(step.OnSuccess)
		}
	} else {
		if step.OnFailure != nil {
			auditStep.OnFailure = executeStep(step.OnFailure)
		}
	}
	return auditStep
}

// evaluateResult if the strings can be parsed into an int, then it will convert the values to ints then it compares the strings with the logicalOperator and returns a status and if needed a description.
func evaluateResult(realValue string, logicalOperator string, expectedValue string) (status string, description string) {

	// Do we compare integers?
	if expectedInt, err := strconv.Atoi(expectedValue); err == nil {
		// We expect a number
		if realValueInt, err := strconv.Atoi(realValue); err == nil {
			// The real value is also a number
			// Then we do the comparison on an int number basis
			switch logicalOperator {
			case "==":
				if realValueInt == expectedInt {
					return success, ""
				} else {
					return fail, ""
				}
			case "!=":
				if realValueInt != expectedInt {
					return success, ""
				} else {
					return fail, ""
				}
			case "<=":
				if realValueInt <= expectedInt {
					return success, ""
				} else {
					return fail, ""
				}
			case ">=":
				if realValueInt >= expectedInt {
					return success, ""
				} else {
					return fail, ""
				}
			default:
				return error, "The logical operator \"" + logicalOperator + "\"  is not Supported"
			}
		}
		if realValue == "" {
			return fail, "Expected int but got string. Trying to compare " + realValue + "(string) " + logicalOperator + " " + expectedValue + "(int). Failed because Type mismatch"
		}
		return error, "Expected int but got string. Trying to compare " + realValue + "(string) " + logicalOperator + " " + expectedValue + "(int). Failed because Type mismatch"
	}

	// We compare the strings
	switch logicalOperator {
	case "==":
		if expectedValue == realValue {
			return success, ""
		} else {
			return fail, ""
		}
	case "!=":
		if expectedValue != realValue {
			return success, ""
		} else {
			return "fail", ""
		}
	case "contains":
		if contains := strings.Contains(realValue, expectedValue); contains {
			return success, ""
		} else {
			return fail, ""
		}
	default:
		return error, "The logical operator \"" + logicalOperator + "\" is not Supported when trying to compare strings"
	}
}
