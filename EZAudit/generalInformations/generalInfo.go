// Package generalInformations implements functions  to get the time, User Info, Device Info, OS Information and if the user is admin.
package generalInformations

import (
	"Just-Go-IT/EZAudit/global"
	"Just-Go-IT/EZAudit/log"
	"os/user"
	"time"
)

// ScanGeneralInfo gets Time, User Info, Device Info, OS Information and returns it as an ReportStruct
func ScanGeneralInfo() (generalInfo global.General) {
	generalInfo.Date = time.Now().String()
	currentUser, err := user.Current()
	if err != nil {
		log.Header("Whilst scanning the current User the following error occurred:\n"+err.Error(), log.Error)
	} else {
		// Information regarding the current User and Device
		generalInfo.User = currentUser.Name + " (ID: " + currentUser.Uid + ")"
		generalInfo.UserName = currentUser.Username
	}
	startOSSpecificScan(&generalInfo)
	return
}
