package generalInformations

import (
	"Just-Go-IT/EZAudit/global"
	"Just-Go-IT/EZAudit/interact"
	"Just-Go-IT/EZAudit/log"
	"os/user"
	"strings"
)

// This function starts the isAdmin() and getOSVersion()
func startOSSpecificScan(general *global.General) {
	general.Admin = isAdmin()
	general.OS = getOSVersion()
}

// Checks if the user has admin rights
func isAdmin() bool {
	u, err := user.Current()
	if err != nil {
		return false
	}
	ids, err := u.GroupIds()
	if err != nil {
		return false
	}
	for i := range ids {
		if ids[i] == "S-1-5-32-544" {
			return true
		}
	}
	return false
}

// This returns detailed Information regarding the Windows Version
func getOSVersion() string {
	combinedOutput, err := interact.Shell("Get-CimInstance -ClassName Win32_OperatingSystem | SELECT *")
	if err != nil {
		log.Header("Whilst trying to get more information about the OS, this error occurred:\n"+err.Error(), log.Error)
	}

	winInfo := make(map[string]string)
	lines := strings.Split(combinedOutput, "\n")
	for _, line := range lines {
		pair := strings.Split(line, " : ")
		if len(pair) > 1 {
			winInfo[strings.TrimSpace(pair[0])] = strings.TrimSpace(pair[1])
		}
	}

	erg := "Product Name: " + winInfo["Caption"] + " | Current Build Version: " + winInfo["BuildNumber"] + " | Version: " + winInfo["Version"]
	return "Windows operating system | " + erg
}
