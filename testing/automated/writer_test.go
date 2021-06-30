package automated

import (
	"Just-Go-IT/EZAudit/artifact"
	"Just-Go-IT/EZAudit/global"
	"Just-Go-IT/EZAudit/writer"
	"os"
	"testing"
)

func Test_WriteResultReport(t *testing.T) {
	var tempstruct global.ReportStruct
	artifact.CreateFolderIfNotExist(global.EZAuditResultDir)
	writer.WriteResultReport(&tempstruct)

	if _, err := os.Stat(global.ResultReportPath); err != nil {
		t.Fail()
	}
	writer.Delete(global.EZAuditResultDir)
}

func Test_CalculateSummary(t *testing.T) {
	testaudits := []global.Audit{
		{
			Name:       "Nummer 1",
			ID:         0,
			AuditSteps: nil,
			Status:     "success",
		},
		{
			Name:       "Nummer 2",
			ID:         1,
			AuditSteps: nil,
			Status:     "fail",
		},
		{
			Name:       "Nummer 3",
			ID:         2,
			AuditSteps: nil,
			Status:     "error",
		},
	}

	testsummary := writer.CalculateSummary(&testaudits)

	if testsummary.Errors != 1 {
		t.Fail()
	}
	if testsummary.Passed != 1 {
		t.Fail()
	}
	if testsummary.Failed != 1 {
		t.Fail()
	}
	if testsummary.PassedPercentage != "33.33%" {
		t.Fail()
	}

}

func Test_RecursiveZip(t *testing.T) {
	artifact.CreateFolderIfNotExist(global.EZAuditResultDir)
	writer.Zip(global.EZAuditResultDir, global.EZAuditResultDir)

	if _, err := os.Stat(global.EZAuditResultDir + ".zip"); err != nil {
		t.Fail()
	}
	writer.Delete(global.EZAuditResultDir)
	os.RemoveAll(global.EZAuditResultDir + ".zip")

}

func Test_CleanUp(t *testing.T) {
	artifact.CreateFolderIfNotExist(global.EZAuditResultDir)

	if _, err := os.Stat(global.EZAuditResultDir); err != nil {
		t.Fail()
	}

	writer.Delete(global.EZAuditResultDir)

	if _, err := os.Stat(global.EZAuditResultDir); err == nil {
		t.Fail()
	}
}
