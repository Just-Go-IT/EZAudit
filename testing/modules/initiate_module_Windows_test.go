package modules

import (
	"Just-Go-IT/EZAudit/global"
	_ "Just-Go-IT/EZAudit/modules/linux"
	_ "Just-Go-IT/EZAudit/modules/windows"
	"testing"
)

func Test_ExecuteAllModulesWindows(t *testing.T) {

	global.ConfigPath = "..\\Testing_modules\\ConfigModuleTestWindows.json"

	/*
		if !strings.Contains(global.ReportStruct.Command[0].Status, "Success") {
			t.Fail()
		}
		if !strings.Contains(global.ReportStruct.Audits[1].Status, "Error") {
			t.Fail()
		}
		if !strings.Contains(global.ReportStruct.Audits[2].Status, "Failed") {
			t.Fail()
		}
		if !strings.Contains(global.ReportStruct.Audits[3].Status, "Failed") {
			t.Fail()
		}
		if !strings.Contains(global.ReportStruct.Audits[4].Status, "Failed") { // Komische Nummer
			t.Fail()
		}
		if !strings.Contains(global.ReportStruct.Audits[5].Status, "Failed") {
			t.Fail()
		}
		if !strings.Contains(global.ReportStruct.Audits[6].Status, "Failed") {
			t.Fail()
		}
		if !strings.Contains(global.ReportStruct.Audits[7].Status, "Failed") {
			t.Fail()
		}
		if !strings.Contains(global.ReportStruct.Audits[8].Status, "Error") { // Weird das normal zu machen
			t.Fail()
		}
		if !strings.Contains(global.ReportStruct.Audits[9].Status, "Error") { // Is error
			t.Fail()
		}
		if !strings.Contains(global.ReportStruct.Audits[10].Status, "Failed") {
			t.Fail()
		}
		if !strings.Contains(global.ReportStruct.Audits[11].Status, "Success") {
			t.Fail()
		}
		if !strings.Contains(global.ReportStruct.Audits[12].Status, "Failed") {
			t.Fail()
		}
		if !strings.Contains(global.ReportStruct.Audits[13].Status, "Failed") {
			t.Fail()
		}
		if !strings.Contains(global.ReportStruct.Audits[14].Status, "Failed") {
			t.Fail()
		}
		if !strings.Contains(global.ReportStruct.Audits[15].Status, "Failed") {
			t.Fail()
		}
		if !strings.Contains(global.ReportStruct.Audits[16].Status, "Failed") {
			t.Fail()
		}
		if !strings.Contains(global.ReportStruct.Audits[17].Status, "Failed") {
			t.Fail()
		}
		if !strings.Contains(global.ReportStruct.Audits[18].Status, "Failed") {
			t.Fail()
		}
		if !strings.Contains(global.ReportStruct.Audits[19].Status, "Failed") {
			t.Fail()
		}
		if !strings.Contains(global.ReportStruct.Audits[20].Status, "Failed") {
			t.Fail()
		}
		if !strings.Contains(global.ReportStruct.Audits[21].Status, "Failed") {
			t.Fail()
		}
		if !strings.Contains(global.ReportStruct.Audits[22].Status, "Failed") {
			t.Fail()
		}
		if !strings.Contains(global.ReportStruct.Audits[23].Status, "Success") {
			t.Fail()
		}
		if !strings.Contains(global.ReportStruct.Audits[24].Status, "Failed") {
			t.Fail()
		}
		if !strings.Contains(global.ReportStruct.Audits[25].Status, "Success") { // Weiss net warum das net tut
			t.Fail()
		}
		if !strings.Contains(global.ReportStruct.Audits[26].Status, "Success") {
				t.Fail()
		}
	*/
	println(global.ArtifactsDir)
}
