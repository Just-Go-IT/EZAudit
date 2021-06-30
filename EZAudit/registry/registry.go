// Package registry provides functions to register modules by name with their struct,
// get module structs by their names and check which os they support and if they need elevation.
package registry

import (
	"Just-Go-IT/EZAudit/flags"
	"Just-Go-IT/EZAudit/global"
	"Just-Go-IT/EZAudit/interact"
	"Just-Go-IT/EZAudit/log"
	"errors"
	"runtime"
	"strconv"
)

// UnidentifiedOS = 1, Linux = 2, Windows = 4, SupportAll = 7
const (
	UnidentifiedOS = 1 << iota
	Linux
	Windows
	SupportAll = UnidentifiedOS | Linux | Windows
)

var reg registry

func init() {
	reg.modules = make(map[string]moduleEntry)
	reg.elevated = interact.IsElevated()
}

// Contains the configOS, the systemOS and the registered modules with the OS they support.
type registry struct {
	modules        map[string]moduleEntry
	configOS       int
	systemOS       int
	elevated       bool
	needsElevation bool
}

type moduleEntry struct {
	module         global.Module
	supportedOS    int
	needsElevation bool
}

// Register the given ModuleName with the given name and the Flag, which OS is supported UnidentifiedOS, Linux, Windows, SupportAll
func Register(name string, module global.Module, needsElevation bool, OSFlag int) {
	reg.modules[name] = moduleEntry{
		module:         module,
		needsElevation: needsElevation,
		supportedOS:    OSFlag,
	}
}

// SetNeedsElevation sets config needs elevation to true
func SetNeedsElevation() {
	reg.needsElevation = true
}

// HasRequiredElevation returns a bool, true means that the required elevation for the config file is fulfilled
func HasRequiredElevation() (requirementsOK bool) {
	var logLevel int
	if reg.needsElevation == false {
		requirementsOK = true
		logLevel = log.Information
	} else {
		requirementsOK = reg.needsElevation == reg.elevated
		if requirementsOK {
			logLevel = log.Information
		} else {
			logLevel = log.Warning
		}
	}
	log.Header("Config requires elevation: "+strconv.FormatBool(reg.needsElevation)+"\n"+
		"Program ran elevated: "+strconv.FormatBool(reg.elevated), logLevel)
	return requirementsOK
}

// GetModule if the module name is registered in the registry it returns the Module. Otherwise it returns nil.
func GetModule(name string) global.Module {
	if ModuleExists(name) {
		return reg.modules[name].module
	}
	return nil
}

// ModuleExists if there is a module in the registry with the passed name it returns true, otherwise false
func ModuleExists(name string) bool {
	_, ok := reg.modules[name]
	return ok
}

// getRecord if the module name is registered in the registry it returns the moduleEntry. Otherwise it returns nil.
func getRecord(name string) *moduleEntry {
	if ModuleExists(name) {
		m := reg.modules[name]
		return &m
	}
	return nil
}

// SetConfigOS sets the passed OS to the registry as the configOS
func SetConfigOS(os string) {
	switch os {
	case "windows":
		reg.configOS = Windows
	case "linux":
		reg.configOS = Linux
	default:
		reg.configOS = UnidentifiedOS
		log.Header("The registry does not recognise the config OS: "+os, log.Error)
	}
}

// setSystemOS sets the passed OS to the registry as the systemOS
func setSystemOS(OSFlag int) {
	reg.systemOS = OSFlag
}

// ModuleSupportsConfigOS  if the module is in the registry and supports the ConfigOS it returns nil. Otherwise it returns an error.
func ModuleSupportsConfigOS(name string) error {
	// Is there a module with this name in the reg?
	if module := getRecord(name); module != nil {
		// Bitwise OR != 0
		if module.supportedOS&reg.configOS != 0 {
			return nil
		}
	}
	return errors.New("the module " + name + " isn´t supported on " + getOS(reg.configOS))
}

// ModuleNeedsElevation returns a bool, true means it needs elevation
func ModuleNeedsElevation(name string) bool {
	if module := getRecord(name); module != nil {
		return module.needsElevation
	}
	return false
}

func systemOSFitsConfigOS() bool {
	return reg.configOS == reg.systemOS
}

// SanityCheck This will check if the ConfigOS matches the SystemOS. If it doesn't match and the force Flag isn't set it will return false.
func SanityCheck() bool {
	// Set the recognized SystemOS
	switch runtime.GOOS {
	case "linux":
		setSystemOS(Linux)
	case "windows":
		setSystemOS(Windows)
	default:
		setSystemOS(UnidentifiedOS)
	}
	if !systemOSFitsConfigOS() {
		if flags.IsSet("force") {
			log.Header("The OS "+runtime.GOOS+" was recognised. But the selected Config-File:"+
				" is for the OS "+getOS(reg.configOS)+"."+
				"\nExecuted anyway because the --force Flag is set.", log.Warning)
		} else {
			log.Header("The OS "+runtime.GOOS+" was recognised. But the selected Config-File: "+
				"is for the OS "+getOS(reg.configOS)+".", log.Error)
			return false
		}
	}
	return true
}

func getOS(os int) string {
	switch os {
	case 1:
		return "unidentified OS"
	case 2:
		return "linux"
	case 4:
		return "windows"
	case 7:
		return "support all OS"
	}
	return "unidentified OS"
}
