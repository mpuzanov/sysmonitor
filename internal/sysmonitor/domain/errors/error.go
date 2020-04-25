package errors

//SysMonitorError для выдачи ошибок системного мониторинга
type SysMonitorError string

func (ee SysMonitorError) Error() string {
	return string(ee)
}

var (
	// ErrReadSystemInfo .
	ErrRunReadSystemInfo = SysMonitorError("Error execute ReadSystemInfo")
	// ErrParserReadSystemInfo .
	ErrParserReadSystemInfo = SysMonitorError("Error parser ReadSystemInfo")
)
