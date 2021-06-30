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
	if u != nil {
		ids, err := u.GroupIds()
		if err != nil {
			return false
		}
		if u.Username == "root" {
			return true
		}
		for i := range ids {
			if strings.Contains(ids[i], "sudo") {
				return true
			}
		}
	}
	return false
}

// This returns detailed Information regarding the Linux Kernel
func getOSVersion() string {
	result, err := interact.Shell("cat /proc/version")
	if err != nil {
		log.Header("Whilst trying to get more information about the OS, this error occurred:\n"+err.Error(), log.Error)
	}
	return strings.TrimSpace(result)
}
