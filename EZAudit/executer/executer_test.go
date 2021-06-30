package executer

import (
	"Just-Go-IT/EZAudit/artifact"
	"Just-Go-IT/EZAudit/global"
	"Just-Go-IT/EZAudit/log"
	_ "Just-Go-IT/EZAudit/modules/linux"
	_ "Just-Go-IT/EZAudit/modules/windows"
	"Just-Go-IT/EZAudit/parser"
	"testing"
)

func Test_EvaluateResult(t *testing.T) {

	tests := []struct {
		number         int
		realValue      string
		comparison     string
		expected       string
		expectedOutput string
	}{
		{1, "10", "==", "10", "success"},
		{2, "1", "!=", "2", "success"},
		{3, "30", "<=", "200", "success"},
		{4, "10", ">=", "5", "success"},

		{5, "123", "==", "1234", "fail"},
		{6, "123", "!=", "123", "fail"},
		{7, "40", "<=", "10", "fail"},
		{8, "123", ">=", "1234", "fail"},
		{9, "123", "<=", "12", "fail"},

		{10, "", ">=", "Buchstabensuppe", "error"},
		{11, "", "!=", "Buchstabensuppe", "success"},
		{12, "", "==", "Buchstabensuppe", "fail"},
		{13, "", "<=", "Buchstabensuppe", "error"},

		{14, "20", "==", "Buchstabensuppe", "fail"},
		{15, "20", "<=", "Buchstabensuppe", "error"},
		{16, "20", ">=", "Buchstabensuppe", "error"},
		{17, "20", "!=", "Buchstabensuppe", "success"},

		{18, "isEqual", "==", "isEqual", "success"},
		{19, "isEqual", "==", "isNotEqual", "fail"},
		{20, "isEqual", "!=", "isNotEqual", "success"},
		{21, "isEqual", "!=", "isEqual", "fail"},

		{22, "Ich bin ein String", "contains", "String", "success"},
		{23, "Ich bin ein String", "contains", "nicht da", "fail"},
	}

	for _, tt := range tests {
		boolOutcome, _ := evaluateResult(tt.realValue, tt.comparison, tt.expected)
		if boolOutcome != tt.expectedOutput {
			t.Fail()
		}
	}
}

func Test_execute(t *testing.T) {
	var config global.ConfigStruct
	global.ConfigPath = "../../testing/automated/testfiles/configExecuteCommands.json"
	parser.ReadConfigFile(global.ConfigPath, &config)
	parser.ParseModules(&config)
	artifact.CreateFolderIfNotExist(global.EZAuditResultDir)
	log.CreateCommandLog(len(config.Commands))

	testaudit := ExecuteAllCommands(&config.Commands)

	if testaudit[0].ID != 0 {
		t.Fail()
	}
	if testaudit[2].Name != "2.3.6.1 (L1) Ensure 'Domain member: Digitally encrypt or sign secure channel data (always)' is set to 'Enabled'" {
		t.Fail()
	}

	if testaudit[0].AuditSteps[0].Status != "fail" {
		t.Fail()
	}
	if testaudit[2].AuditSteps[0].Status != "success" {
		t.Fail()
	}
	println("waiting for concurrency")
}
