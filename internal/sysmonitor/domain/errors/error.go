package errors

//SysMonitorError для выдачи ошибок системного мониторинга
type SysMonitorError string

func (ee SysMonitorError) Error() string {
	return string(ee)
}

var (
	// ErrRunReadInfoSystem .
	ErrRunReadInfoSystem = SysMonitorError("Error execute ReadInfoSystem")
	// ErrParserReadInfoSystem .
	ErrParserReadInfoSystem = SysMonitorError("Error parser ReadInfoSystem")

	// ErrRunReadInfoCPU .
	ErrRunReadInfoCPU = SysMonitorError("Error execute ReadInfoCPU")
	// ErrParserReadInfoCPU .
	ErrParserReadInfoCPU = SysMonitorError("Error parser ReadInfoCPU")

	// ErrRunLoadDiskDevice .
	ErrRunLoadDiskDevice = SysMonitorError("Error execute RunLoadDiskDevice")
	// ErrParserLoadDiskDevice .
	ErrParserLoadDiskDevice = SysMonitorError("Error parser ReadLoadDiskDevice")

	// ErrRunLoadDiskFS .
	ErrRunLoadDiskFS = SysMonitorError("Error execute RunLoadDiskFS")
	// ErrParserLoadDiskFS .
	ErrParserLoadDiskFS = SysMonitorError("Error parser ReadLoadDiskFS")

	// ErrRunLoadDiskFSInode .
	ErrRunLoadDiskFSInode = SysMonitorError("Error execute RunLoadDiskFSInode")
	// ErrParserLoadDiskFSInode .
	ErrParserLoadDiskFSInode = SysMonitorError("Error parser ReadLoadDiskFSInode")

	// ErrRunDeviceNet .
	ErrRunDeviceNet = SysMonitorError("Error execute ErrRunDeviceNet")
	// ErrParserDeviceNet .
	ErrParserDeviceNet = SysMonitorError("Error parser ErrParserDeviceNet")

	// ErrRunNetworkStatistics .
	ErrRunNetworkStatistics = SysMonitorError("Error execute ErrRunNetworkStatistics")
	// ErrParserNetworkStatistics .
	ErrParserNetworkStatistics = SysMonitorError("Error parser ErrParserNetworkStatistics")
)
