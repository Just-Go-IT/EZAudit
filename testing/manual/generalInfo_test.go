package manual

import (
	"Just-Go-IT/EZAudit/generalInformations"
	"reflect"
	"testing"
)

func Test_ScanGeneralInfo(t *testing.T) {
	// nur ausführbar auf NB-Win10
	report := generalInformations.ScanGeneralInfo()

	if !reflect.DeepEqual(report.User, " (ID: S-1-5-21-2068594027-868955951-3733561519-1001)") {
		t.Fail()
	}
	if report.UserName != "NB-WIN10\\Nicolas Biundo" {
		t.Fail()
	}

	// OS Specific Scan Entrys
	if report.Admin != true {
		t.Fail()
	}
	if report.OS != "Windows operating system | Product Name: Microsoft Windows 10 Pro | Current Build Version: 19042 | Version: 10.0.19042" {
		t.Fail()
	}

}
