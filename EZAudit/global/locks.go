package global

import "sync"

var (
	LoggerLock          *sync.Mutex
	SeceditLock         *sync.Mutex
	AuditpolLock        *sync.Mutex
	IsInstalledLock     *sync.Mutex
	WindowsServicesLock *sync.Mutex
	PipeStorageLock     *sync.Mutex
	CommandMapLock      *sync.Mutex
)

func init() {
	LoggerLock = &sync.Mutex{}
	SeceditLock = &sync.Mutex{}
	AuditpolLock = &sync.Mutex{}
	IsInstalledLock = &sync.Mutex{}
	WindowsServicesLock = &sync.Mutex{}
	PipeStorageLock = &sync.Mutex{}
	CommandMapLock = &sync.Mutex{}
}
