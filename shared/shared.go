package shared

const (
	// InformationMessage is the hint for informational messages.
	InformationMessage = "INFO"
	// WarningMessage is the hint for warning messages.
	WarningMessage = "WARN"
	// CriticalMessage is the hint for critical level messages.
	CriticalMessage = "CRIT"
	// DebugMessage is the hint for debug level messages.
	DebugMessage = "DEBUG"
)

// AuditSender is not able to drop messages, it will therefore slow down your
// application in order to ship logs. This can have undesirable effects on your application,
// however if logs are more important than service then this is the only option.
func AuditSender(msg string, pipe chan string) {
	pipe <- msg
}

// BestEffortSender will send messages to the printer until the buffer is full. Once the buffer
// is full messages will spill and be dropped. Each message dropped will in increase the dropped
// message counter. This is mode is useful when you decide that service is more important than
// log shipping. Most time users will want this option even though they may not have thought
// much about it. It is therefore the default option.
func BestEffortSender(msg string, pipe chan string, droppedFunc func()) {
	select {
	case pipe <- msg:
		return
	default:
		droppedFunc()
	}
}
