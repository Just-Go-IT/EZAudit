// Package parser implements all functions to read and parse a config json file to the internal structs
package parser

import (
	"Just-Go-IT/EZAudit/flags"
	"Just-Go-IT/EZAudit/global"
	"Just-Go-IT/EZAudit/log"
	"Just-Go-IT/EZAudit/modules"
	_ "Just-Go-IT/EZAudit/modules/linux"   //register the modules in the registry
	_ "Just-Go-IT/EZAudit/modules/windows" //register the modules in the registry
	"Just-Go-IT/EZAudit/registry"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
)

// CheckConfigExists If the Config File is available it returns, otherwise the function asks for a new path and calls itself again til the path is correct
func CheckConfigExists(configPath *string) bool {
	// Is the config.json file available?
	if _, err := os.Stat(*configPath); err != nil {
		fmt.Println("The expected config file \"" + *configPath + "\" was not found. " +
			"Please enter the path to the config file (for example:" +
			"/files/config.json)")

		path, _ := os.Getwd()

		fmt.Print(path + ">")
		_, scanErr := fmt.Scan(configPath)

		if scanErr != nil {
			os.Stderr.WriteString("Whilst scanning the UserInput for the ConfigPath this error occurred:\n" + err.Error())
		}

		return false
	}
	return true
}

// ReadConfigFile reads the Config-File and parses it to a global.ConfigStrut. Returns false if there was an error.
func ReadConfigFile(configPath string, config *global.ConfigStruct) (ok bool) {
	// Open the Config File
	configFile, err := os.Open(configPath)
	if err != nil {
		log.Header("Whilst trying to open "+configPath+" this Error occurred:\n"+err.Error(), log.Fatal)
		return false
	} else {
		defer func(configFile *os.File) {
			errClose := configFile.Close()
			if errClose != nil {
				log.Header("Whilst trying to close "+configPath+" this Error occurred:\n"+errClose.Error(), log.Error)
			}
		}(configFile)
	}

	// Read the opened jsonFile as a byte array
	byteConfigFile, errRead := ioutil.ReadAll(configFile)
	if errRead != nil {
		log.Header("Whilst trying to read "+configPath+" this Error occurred:\n"+errRead.Error(), log.Fatal)
		return false
	}

	// Parse the opened jsonFile to the internal global.ConfigStruct
	errJson := json.Unmarshal(byteConfigFile, config)
	if errJson != nil {
		log.Header("Whilst parsing "+configPath+" this Error occurred:\n"+errJson.Error(), log.Fatal)
		return false
	}
	return true
}

// ParseModules Parses functional structs for every Step. Creates Modules and saves it to the internal Struct
func ParseModules(config *global.ConfigStruct) (errCounter int) {
	// The dry run relates to the configOS
	registry.SetConfigOS(config.OS)

	for i := range config.Commands { // Iterate through the Commands
		config.Commands[i].Index = i
		for j := range config.Commands[i].Steps { // Iterate through the Command Steps
			if j < len(config.Commands[i].Steps)-1 {
				config.Commands[i].Steps[j].GetNext = &config.Commands[i].Steps[j+1] // Save the next Step to the internal Step struct for Piping
			}
			// Every Step has its own Path with all Information that it needs
			config.Commands[i].Steps[j].Path.CommandName = config.Commands[i].Name
			config.Commands[i].Steps[j].Path.CommandIndex = i
			config.Commands[i].Steps[j].Path.StepIndex = j

			errCounter += createModule(&config.Commands[i].Steps[j]) // Creates a Module and saves it to Step struct
		}
	}
	return errCounter
}

// createModule validates moduleName and the parameters. Checks if the all requirements for the module are given, and logs all the information.
// If everything is OK- create a functional module and save it to Step struct, else create a BrokenModule.
func createModule(step *global.Step) (errCounter int) {
	// Is the Module registered in the registry? (does it exist?)
	if !registry.ModuleExists(step.ModuleName) {
		errCounter++
		log.DryRunError("there is no module in the registry called "+step.ModuleName, step.Path)
		step.Module = modules.BrokenModule{Reason: "there is no module in the registry called " + step.ModuleName}
	} else {
		// Does the Module Support the Config OS?
		err := registry.ModuleSupportsConfigOS(step.ModuleName)
		if err != nil {
			// Module does not support the OS
			errCounter++
			log.DryRunError(err.Error(), step.Path)
		}
		if err != nil && !flags.IsSet("force") {
			step.Module = modules.BrokenModule{Reason: err.Error()}
		} else {
			// Module is in the registry and supports the config OS
			if step.Parameters == nil {
				errCounter++
				log.DryRunError("the object parameter is missing in the step.", step.Path)
			}
			var errM error
			step.Module, errM = registry.GetModule(step.ModuleName).New(step.Parameters)
			// Could the parameters be parsed?
			if errM != nil {
				errCounter++
				log.DryRunError(errM.Error(), step.Path)
				step.Module = modules.BrokenModule{Reason: "This Error occurred whilst parsing the parameters: " + errM.Error()}
			} else {
				if step.NeedsElevation || registry.ModuleNeedsElevation(step.ModuleName) {
					step.NeedsElevation = true
					registry.SetNeedsElevation()
					log.DryRunError("The module"+step.ModuleName+" will need elevation", step.Path)
					errCounter++
				}
			}
		}
	}

	// Parse Regex for censoring
	for _, reg := range step.Censor {
		tmp, err := regexp.Compile(reg)
		if err != nil {
			errCounter++
			log.DryRunError("the regex syntax is incorrect. While compiling, this error occurred\n"+err.Error(), step.Path)
		} else {
			step.Regex = append(step.Regex, *tmp)
		}
	}

	// If there are Other Steps inside the Step, continue going deeper
	if step.OnSuccess != nil { // Recursion if OnSucces
		step.OnSuccess.Path = step.Path
		step.OnSuccess.Path.ExactPath += "_OnSuccess"
		errCounter += createModule(step.OnSuccess)
	}
	if step.OnFailure != nil { // Recursion if OnFailure
		step.OnFailure.Path = step.Path
		step.OnFailure.Path.ExactPath += "_OnFailure"
		errCounter += createModule(step.OnFailure)
	}

	return errCounter
}
