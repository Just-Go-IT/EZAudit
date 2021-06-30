package main

import (
	"Just-Go-IT/EZAudit/artifact"
	"Just-Go-IT/EZAudit/executer"
	"Just-Go-IT/EZAudit/flags"
	"Just-Go-IT/EZAudit/generalInformations"
	"Just-Go-IT/EZAudit/global"
	"Just-Go-IT/EZAudit/log"
	"Just-Go-IT/EZAudit/parser"
	"Just-Go-IT/EZAudit/registry"
	"Just-Go-IT/EZAudit/writer"
	"fmt"
	"strconv"
	"time"
)

func main() {
	start := time.Now()

	flags.ParseFlags()

	log.ActiveFlags()

	// Prints the --help menu in the console and exits
	if flags.IsSet("help") {
		flags.PrintCommandLineHelp()
		return
	}

	fmt.Println("Starting EZAudit ( https://github.com/Just-Go-IT/ERNW-Audit-Framework ) at " + start.Format("2006-01-02 15:04"))

	if !flags.IsSet("noZip") {
		defer func() {
			// Zip the result folder
			err := writer.Zip(global.EZAuditResultDir, global.EZAuditResultDir)
			if err != nil {
				println("Whilst zipping all the files the following Error occurred:\n" + err.Error())
			}
			// Clean up the unzipped folders
			writer.Delete(global.EZAuditResultDir)
		}()
	}

	// Always write a debug.log at the end of the program
	defer log.Write(global.DebugLogPath)

	// Was a config file passed in the arguments? If yes, use the passed config file path
	if flags.IsSet("config") {
		f := flags.GetFlagByName("config")
		global.ConfigPath = f.GetValue()
	}

	// Ensures that the config path is correct after this point
	for !parser.CheckConfigExists(&global.ConfigPath) {
	}

	artifact.CreateFolderIfNotExist(global.EZAuditResultDir)

	var config global.ConfigStruct

	// Reads the Config File and unmarshal it
	if !parser.ReadConfigFile(global.ConfigPath, &config) {
		fmt.Println("The config file couldn't be read. Check the debug.log")
		return
	}

	// The log needs the size of the commands to create an array
	log.CreateCommandLog(len(config.Commands))

	// Set the verbosity Level for the logger but the flags have higher priority
	if flags.IsSet("verbosity") {
		v, _ := strconv.Atoi(flags.GetFlagByName("verbosity").GetValue())
		log.SetVerbosityLevel(v)
	} else {
		log.SetVerbosityLevel(config.Verbosity)
	}

	// Parses the steps to modules and counts the failures
	errCounter := parser.ParseModules(&config)

	if flags.IsSet("dryRun") {
		log.DryRunResult(errCounter)
		fmt.Println("DryRun finished.")
		if errCounter == 0 {
			fmt.Println("The config file is valid.")
		} else {
			fmt.Println(fmt.Sprintf("There are %v error(s) in the config file."+
				"For more information, check the debug.log", errCounter))
		}
		return
	}

	// Copies the config file to the result folder
	defer func() {
		err := writer.CopyFile(global.ConfigPath, global.ConfigResultPath)
		if err != nil {
			log.Footer(err.Error(), log.Error)
		}
	}()

	artifact.SetUp(config.Censor, config.MaxOutputLength)

	// Check if the config file OS matches the running OS
	if !registry.SanityCheck() {
		fmt.Println("Didn't pass the sanity check, check the debug.log")
		return
	}

	if !registry.HasRequiredElevation() && !flags.IsSet("force") {
		fmt.Println("This Configuration needs elevation, please restart the program with elevated rights (sudo), or use --force")
		return
	}

	var report global.ReportStruct

	// Starts an OS-specific Scan for the Header
	report.General = generalInformations.ScanGeneralInfo()

	// Processing Commands to Audits
	report.Audits = executer.ExecuteAllCommands(&config.Commands)

	// Writes the Summary
	report.Summary = writer.CalculateSummary(&report.Audits)

	log.Summary(report.Summary)

	// Tracks the time that passed
	elapsedTime := fmt.Sprintf("Elapsed time %0.2fs", time.Since(start).Seconds())
	report.General.Runtime = elapsedTime

	// Write report and zip all files
	writer.WriteResultReport(&report)

	// Prints the Summary
	fmt.Println(fmt.Sprintf("\nSummary:\n"+
		"------------------------------------------------------------------------------------------------------------------------\n"+
		"PassedPercentage: %v\n"+
		"\tPassed: %v\n"+
		"\tFailed: %v\n"+
		"\tErrors: %v\n"+
		"------------------------------------------------------------------------------------------------------------------------\n"+
		"Elapsed time: %v"+
		"\nThe results are in the File %v", report.Summary.PassedPercentage, report.Summary.Passed, report.Summary.Failed,
		report.Summary.Errors, elapsedTime, global.EZAuditResultDir))
}
