package automated

import (
	"Just-Go-IT/EZAudit/global"
	log2 "Just-Go-IT/EZAudit/log"
	"Just-Go-IT/EZAudit/parser"
	"io/ioutil"
	"log"
	"os"
	"testing"
)

func Test_ParseConfigFile(t *testing.T) {
	global.ConfigPath = "testfiles\\goodConfig.json"

	var tempstruct global.ConfigStruct

	parser.ReadConfigFile(global.ConfigPath, &tempstruct)

	if tempstruct.OS != "windows" {
		t.Fail()
	}
	if tempstruct.Commands[0].Steps[0].ModuleName != "secedit" {
		t.Fail()
	}
	if tempstruct.Commands[1].Steps[0].ModuleName != "secedit" {
		t.Fail()
	}
	if tempstruct.Commands[0].Steps[0].Comparison != "==" {
		t.Fail()
	}
	if tempstruct.Commands[0].Steps[0].ExpectedValue != "24" {
		t.Fail()
	}
}

func Test_CheckConfigPathExists(t *testing.T) {

	tests := []struct {
		name           string
		configfilepath string
		userinput      []byte
	}{
		{"conf in default path", "../../configExamples/windowsExample.json", []byte("NoCAll")},
		{"conf not default, correct input", "wrong/config.json", []byte("../../configExamples/windowsExample.json")},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmpfile, err := ioutil.TempFile("", "example")
			if err != nil {
				log.Fatal(err)
			}

			defer os.Remove(tmpfile.Name()) // clean up

			if _, err1 := tmpfile.Write(tt.userinput); err1 != nil {
				log.Fatal(err)
			}

			if _, err2 := tmpfile.Seek(0, 0); err2 != nil {
				log.Fatal(err)
			}

			oldStdin := os.Stdin
			os.Stdin = tmpfile
			global.ConfigPath = tt.configfilepath
			parser.CheckConfigExists(&global.ConfigPath)
			if err != nil {
				t.Fail()
				t.Fatal(err)
			}

			os.Stdin = oldStdin
		})
	}
}

func Test_ParseModules(t *testing.T) {
	// calls CreateModules and SanitizeDirName

	tests := []struct {
		configfilepath string
		errcounter     int
	}{
		{"testfiles\\configExecuteCommands.json", 0},
		{"testfiles\\wrongModulesConfig.json", 2},
	}

	var tempstruct global.ConfigStruct

	for _, tt := range tests {
		parser.ReadConfigFile(tt.configfilepath, &tempstruct)
		log2.CreateCommandLog(len(tempstruct.Commands))
		counter := parser.ParseModules(&tempstruct)
		if tt.errcounter != counter {
			t.Fail()
		}
	}

	if tempstruct.Commands[0].Steps[0].Path.CommandName != "1.1.1 (L1) Ensure 'Enforce password history' is set to '24 or more password(s)" {
		t.Fail()
	}

}
