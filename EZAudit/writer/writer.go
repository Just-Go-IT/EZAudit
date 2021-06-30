// Package writer provides functions to exports the result report.
package writer

import (
	"Just-Go-IT/EZAudit/global"
	"Just-Go-IT/EZAudit/log"
	"archive/zip"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

// CalculateSummary counts how many audits were successful, failed or errored. Also calculates the how many audits passed in percentage.
func CalculateSummary(audits *[]global.Audit) (summary global.Summary) {
	var success, failed, err int
	for _, audit := range *audits {
		switch audit.Status {
		case "success":
			success++
		case "fail":
			failed++
		case "error":
			err++
		}
	}
	// If divided by zero.
	if success+failed+err <= 0 {
		summary.PassedPercentage = "0.00%"
	} else {
		summary.PassedPercentage = strconv.FormatFloat(float64(success)/float64(success+failed+err)*100, 'f', 2, 64) + "%"
	}
	summary.Passed = success
	summary.Failed = failed
	summary.Errors = err
	return
}

// WriteResultReport Writes all the results in a resultReport.json File.
func WriteResultReport(report *global.ReportStruct) {
	file, _ := json.MarshalIndent(report, "", " ")

	// adjust encoding for "<" and ">"
	file = bytes.Replace(file, []byte("\\u003c"), []byte("<"), -1)
	file = bytes.Replace(file, []byte("\\u003e"), []byte(">"), -1)

	err := ioutil.WriteFile(filepath.FromSlash(global.ResultReportPath), file, 0644)
	if err != nil {
		log.Footer("Whilst writing the file "+global.ResultReportPath+" the following Error occurred:\n"+err.Error(), log.Error)
		fmt.Println(err.Error())
	}
}

// CopyFile copies a file from a source to a destination and returns an error if there was an error, otherwise nil
func CopyFile(source, destination string) error {
	input, err := ioutil.ReadFile(source)
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(destination, input, global.FilePermissionValue)
	if err != nil {
		return err
	}
	return nil
}

// Zip Creates a zipped copy of the the source and saves it in destination.
func Zip(destination string, source string) error {
	// Creates the zip File
	destinationFile, err := os.Create(destination + ".zip")
	if err != nil {
		return err
	}
	// Writes it into the created .zip File
	zipWriter := zip.NewWriter(destinationFile)
	// Defer closes the zipWriter
	defer func(myZip *zip.Writer) {
		errZip := myZip.Close()
		if errZip != nil {

		}
	}(zipWriter)

	// Walk the filepath through the source to copy every file and folder to the destination zip file
	err = filepath.Walk(source, func(filePath string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}
		if err != nil {
			return err
		}
		relPath := strings.TrimPrefix(filePath, source+"\\")
		zipFile, err := zipWriter.Create(relPath)
		if err != nil {
			return err
		}
		fsFile, err := os.Open(filePath)
		if err != nil {
			return err
		}
		defer fsFile.Close()
		_, err = io.Copy(zipFile, fsFile)
		if err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		return err
	}
	err = zipWriter.Close()
	if err != nil {
		return err
	}
	return nil
}

// Delete removes all files and folders in the specified path
func Delete(path string) {
	err := os.RemoveAll(path)
	if err != nil {
		println("While trying to remove " + path + " the following error occurred:\n" + err.Error())
	}
}
