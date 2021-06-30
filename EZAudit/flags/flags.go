// Package flags implements functions for parsing passed program arguments.
package flags

import (
	"errors"
	"fmt"
	"os"
	"strings"
)

type flag struct {
	name         string
	syntax       []string
	description  string
	isSet        bool
	deliverValue bool
	value        string
}

var flags = []flag{
	{
		name:        "dryRun",
		syntax:      []string{"-d", "--dryRun"},
		description: "the DryRun Mode validates the Config File. Checks if JSON format is correct,\n\tif Modules have correct Parameters, and if the OS is matching.",
	},
	{
		name:        "force",
		syntax:      []string{"-f", "--force"},
		description: "the force Mode forces the program to execute the config file.",
	},
	{
		name:         "config",
		syntax:       []string{"-c=", "--config="},
		description:  "pass the path of the config file. For example --config=/Desktop/config2.json",
		deliverValue: true,
	},
	{
		name:         "verbosity",
		syntax:       []string{"-v=", "--verbosity="},
		description:  "sets the verbosity for the logging.\n\t1 - Fatal\n\t2 - Error\n\t3 - Warning\n\t4 - Information(Default)\n\t5 - Debug",
		deliverValue: true,
	},
	{
		name:        "noZip",
		syntax:      []string{"-nz", "--noZip"},
		description: "don't zip the result folder",
	},
	{
		name:        "help",
		syntax:      []string{"-h", "--help"},
		description: "gives command Line help",
	},
}

// ParseFlags sets all flags in the os.Args if they are supported, else if the syntax is wrong it transfers the user to the --help menu
func ParseFlags() {
	for i, userFlag := range os.Args {
		// Skips the first flag, because it is always in the execution path
		if i == 0 {
			continue
		}
		// If there are more flags, parse them
		err := parseFlag(userFlag)
		if err != nil {
			// Transfers the user to the help menu
			fmt.Println(err.Error())
			parseFlag("--help")
		}
	}
}

// parseFlag checks if flag is supported and the syntax is correct. sets flag if supported and returns true, else returns false.
func parseFlag(userFlag string) (err error) {
	// If the flag has a value, call function to handle it accordingly
	// The the flag has no value, try to match with supported flags
	if strings.Contains(userFlag, "=") {
		return handleFlagWithValue(userFlag)
	}
	if f := getFlagBySyntax(userFlag); f != nil {
		f.set(true)
		return nil
	}
	// No match means that the flag does not exist, or the syntax is incorrect
	return errors.New("Wrong Syntax. The flag " + userFlag + " does not exist")
}

// handleFlagWithValue if the userFlag contains a value (behind an equals "=" operator). Handles the flag and sets the given value to the specific inner flags.
// Returns nil if not supported or wrong syntax
func handleFlagWithValue(userFlag string) error {
	s := strings.Split(userFlag, "=")
	s[0] += "="

	if s[0] == "-v=" || s[0] == "--verbosity=" {
		if s[1] < "1" || s[1] > "5" {
			return errors.New("Wrong Syntax. The flag " + s[0] + " does not support the value " + s[1])
		}
	}

	if f := getFlagBySyntax(s[0]); f != nil {
		f.value = s[1]
		f.set(true)
		return nil
	} else {
		return errors.New("Wrong Syntax. The flag " + userFlag + " does not exist")
	}
}

// getFlagBySyntax iterates trough inner flags and matches the syntax. If match found, returns the flag, otherwise returns nil.
func getFlagBySyntax(userFlag string) *flag {
	for i, f := range flags {
		for _, syntax := range f.syntax {
			if strings.ToLower(userFlag) == strings.ToLower(syntax) {
				return &flags[i]
			}
		}
	}
	return nil
}

// GetFlagByName iterates trough inner flags and matches the name. If match found, returns the flag, otherwise returns nil.
func GetFlagByName(flagName string) *flag {
	for _, f := range flags {
		if f.name == flagName {
			return &f
		}
	}
	return nil
}

// GetAllSetFlags returns all set Flags
func GetAllSetFlags() (setFlags []flag) {
	for _, f := range flags {
		if f.isSet {
			setFlags = append(setFlags, f)
		}
	}
	return
}

// IsSet returns true if the flag is isSet, else false
func IsSet(flagName string) bool {
	if f := GetFlagByName(flagName); f != nil {
		return f.isSet
	}
	return false
}

// PrintCommandLineHelp prints the command line help
func PrintCommandLineHelp() {
	fmt.Println("Usage of EZAudit:")
	for _, f := range flags {
		fmt.Println("Syntax: " + strings.ReplaceAll(fmt.Sprintf("%v", f.syntax), " ", ", "))
		fmt.Println("\t" + f.description)
	}
}

// set changes the flag.isSet parameter
func (f *flag) set(b bool) {
	f.isSet = b
}

// GetValue returns the value of the flag
func (f flag) GetValue() string {
	return f.value
}

// GetName returns the name of the flag
func (f flag) GetName() string {
	return f.name
}
