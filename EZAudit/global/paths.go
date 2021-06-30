package global

import (
	"os"
	"path/filepath"
	"time"
)

var (
	FilePermissionValue os.FileMode = 0777

	//Paths
	RunningDir          string
	EZAuditResultDir    string
	ConfigPath          string
	ConfigResultPath    string
	DebugLogPath        string
	ArtifactsDir        string
	AuditDir            string
	MainResourcesDir    string
	ResultReportPath    string
	SeceditDir          string
	SeceditPath         string
	AuditpolDir         string
	AuditpolPath        string
	IsInstalledDir      string
	IsInstalledPath     string
	WindowsServicesDir  string
	WindowsServicesPath string
)

func init() {
	ConfigPath = "config.json"
	RunningDir = filepath.Dir(os.Args[0])
	EZAuditResultDir = filepath.Join(RunningDir, "EZAuditResult"+time.Now().Format("_02-Jan-2006_15-04-05"))
	DebugLogPath = filepath.Join(EZAuditResultDir, "debug.log")
	ConfigResultPath = filepath.Join(EZAuditResultDir, "config.json")
	ResultReportPath = filepath.Join(EZAuditResultDir, "resultReport.json")
	ArtifactsDir = filepath.Join(EZAuditResultDir, "artifacts")
	AuditDir = filepath.Join(ArtifactsDir, "audit")
	MainResourcesDir = filepath.Join(ArtifactsDir, "mainResources")
	SeceditDir = filepath.Join(MainResourcesDir, "secedit")
	SeceditPath = filepath.Join(SeceditDir, "secedit.txt")
	IsInstalledDir = filepath.Join(MainResourcesDir, "isInstalled")
	IsInstalledPath = filepath.Join(IsInstalledDir, "isInstalled.txt")
	WindowsServicesDir = filepath.Join(MainResourcesDir, "windowsServices")
	WindowsServicesPath = filepath.Join(WindowsServicesDir, "windowsServices.txt")
	AuditpolDir = filepath.Join(MainResourcesDir, "auditpol")
	AuditpolPath = filepath.Join(AuditpolDir, "auditpol.txt")
}
