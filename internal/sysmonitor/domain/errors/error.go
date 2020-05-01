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
)
