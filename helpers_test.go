package loggos

func shutdownCurrentLoggers() {
	if DefaultJSONLogger != nil || DefaultLineLogger != nil {
		<-Flush()
	}
	DefaultJSONLogger = nil
	DefaultLineLogger = nil
}
