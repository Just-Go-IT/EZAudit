package automated

import (
	"Just-Go-IT/EZAudit/artifact"
	"Just-Go-IT/EZAudit/global"
	"Just-Go-IT/EZAudit/log"
	"Just-Go-IT/EZAudit/writer"
	"os"
	"testing"
)

func Test_CreateFolderIfNotExist(t *testing.T) {
	writer.Delete(global.EZAuditResultDir)
	if _, err := os.Stat(global.EZAuditResultDir); err == nil {
		t.Fail()
	}
	artifact.CreateFolderIfNotExist(global.EZAuditResultDir)
	if _, err := os.Stat(global.EZAuditResultDir); err != nil {
		t.Fail()
	}
	writer.Delete(global.EZAuditResultDir)
}

func Test_SaveString(t *testing.T) {
	artifact.CreateFolderIfNotExist(global.EZAuditResultDir)
	log.CreateCommandLog(5)
	testString := "Das ist ein Teststring um zu schauen ob der Atifact Collecter ordnungsgemäß funktioniert"

	teststep := global.Step{
		ModuleName:       "",
		AllowFailure:     false,
		DontSaveArtifact: false,
		UsePipe:          false,
		Parameters:       nil,
		Comparison:       "",
		ExpectedValue:    "",
		Censor:           nil,
		OnSuccess:        nil,
		OnFailure:        nil,
		NeedsElevation:   false,
		Path: global.Path{
			CommandName:  "Testcommand1",
			CommandIndex: 0,
			StepIndex:    0,
			ExactPath:    "Test_OnSuccess",
		},
		Module:  nil,
		Regex:   nil,
		GetNext: nil,
	}

	artifact.SaveString(testString, teststep)
	if _, err := os.Stat(global.AuditDir + "\\_000_Testcommand1\\0Test_OnSuccess.txt"); err != nil {
		t.Fail()
	}
	if data, err := os.ReadFile(global.AuditDir + "\\_000_Testcommand1\\0Test_OnSuccess.txt"); err == nil {
		if string(data[:]) != "Das ist ein Teststring um zu schauen ob der Atifact Collecter ordnungsgemäß funktioniert" {
			t.Fail()
		}
	}

	artifact.SaveString(testString, teststep)
	if _, err := os.Stat(global.AuditDir + "\\_000_Testcommand1\\0Test_OnSuccess_1.txt"); err != nil {
		t.Fail()
	}

	teststep.Path.ExactPath = "DirectString"
	artifact.SaveString("testString", teststep)
	if data, err := os.ReadFile(global.AuditDir + "\\_000_Testcommand1\\0DirectString.txt"); err == nil {
		if string(data[:]) != "testString" {
			t.Fail()
		}
	}
	writer.Delete(global.EZAuditResultDir)
}

func Test_SaveFile(t *testing.T) {
	artifact.CreateFolderIfNotExist(global.EZAuditResultDir)
	log.CreateCommandLog(5)
	path := "testfiles\\teststring.txt"

	teststep := global.Step{
		ModuleName:       "",
		AllowFailure:     false,
		DontSaveArtifact: false,
		UsePipe:          false,
		Parameters:       nil,
		Comparison:       "",
		ExpectedValue:    "",
		Censor:           nil,
		OnSuccess:        nil,
		OnFailure:        nil,
		NeedsElevation:   false,
		Path: global.Path{
			CommandName:  "Testcommand1",
			CommandIndex: 0,
			StepIndex:    0,
			ExactPath:    "Test_OnSuccess",
		},
		Module:  nil,
		Regex:   nil,
		GetNext: nil,
	}

	artifact.SaveFile(path, teststep)
	if _, err := os.Stat(global.AuditDir + "\\_000_Testcommand1\\0Test_OnSuccess.txt"); err != nil {
		t.Fail()
	}
	artifact.SaveFile(path, teststep)
	if _, err := os.Stat(global.AuditDir + "\\_000_Testcommand1\\0Test_OnSuccess_1.txt"); err != nil {
		t.Fail()
	}
	if data, err := os.ReadFile(global.AuditDir + "\\_000_Testcommand1\\0Test_OnSuccess.txt"); err == nil {
		if string(data[:]) != "I'm your testfile!" {
			t.Fail()
		}
	}
	writer.Delete(global.EZAuditResultDir)
}

func Test_CensorValueGlobal(t *testing.T) {
	regexExp := []string{"password:.*"}
	artifact.SetUp(regexExp, -1)
	censorstring := "password: pa$$w0rd"

	output := artifact.CensorGlobal(censorstring)

	if output != "--censored--" {
		println(output)
		t.Fail()
	}
}
