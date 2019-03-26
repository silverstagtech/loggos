package loggos

import "github.com/silverstagtech/loggos/jsonmessage"

// The below functions are all shortcuts that relate to the default logger.
// Look at the function comments for the JSON Logger type for details on what they do.

// SendJSON is used to send a JSON message to the default JSON logger.
func SendJSON(msg *jsonmessage.JSONMessage) {
	startdefaultJSONLogger()
	DefaultJSONLogger.Send(msg)
}

// JSONLoggerEnableDebugLogging Starts the default JSON logger if not already started then
// enabled debug logging.
func JSONLoggerEnableDebugLogging(toggle bool) {
	startdefaultJSONLogger()
	DefaultJSONLogger.EnableDebugLogging(toggle)
}

// JSONLoggerEnablePrettyPrint starts the default JSON logger if not already started then
// enabled debug logging
func JSONLoggerEnablePrettyPrint(toggle bool) {
	startdefaultJSONLogger()
	DefaultJSONLogger.EnablePrettyPrint(toggle)
}
