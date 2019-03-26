package loggos

// The below functions are all shortcuts that relate to the default logger.
// Look at the function comments for the Line Logger type for details on what they do.

func LineLoggerEnableDebugLogging(toggle bool) {
	startdefaultLineLogger()
	DefaultLineLogger.EnableDebugLogging(toggle)
}

func Infoln(msg ...interface{}) {
	// Start the line logger if needed.
	startdefaultLineLogger()
	DefaultLineLogger.Infoln(msg...)

}
func Warnln(msg ...interface{}) {
	// Start the line logger if needed.
	startdefaultLineLogger()
	DefaultLineLogger.Warnln(msg...)

}
func Critln(msg ...interface{}) {
	// Start the line logger if needed.
	startdefaultLineLogger()
	DefaultLineLogger.Critln(msg...)

}
func Debugln(msg ...interface{}) {
	// Start the line logger if needed.
	startdefaultLineLogger()
	DefaultLineLogger.Debugln(msg...)

}
func Infof(format string, vars ...interface{}) {
	// Start the line logger if needed.
	startdefaultLineLogger()
	DefaultLineLogger.Infof(format, vars...)
}
func Warnf(format string, vars ...interface{}) {
	// Start the line logger if needed.
	startdefaultLineLogger()
	DefaultLineLogger.Warnf(format, vars...)
}
func Critf(format string, vars ...interface{}) {
	// Start the line logger if needed.
	startdefaultLineLogger()
	DefaultLineLogger.Critf(format, vars...)
}
func Debugf(format string, vars ...interface{}) {
	// Start the line logger if needed.
	startdefaultLineLogger()
	DefaultLineLogger.Debugf(format, vars...)
}
