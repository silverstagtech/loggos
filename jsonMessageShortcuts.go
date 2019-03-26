package loggos

import "github.com/silverstagtech/loggos/jsonmessage"

// JSONInfoln is a shortcut function that will give you a JSON Message that is populated with time, level
// and your message ready to be shipped. It can still be decorated with more keys if needed.
func JSONInfoln(msg ...interface{}) *jsonmessage.JSONMessage {
	jmsg := jsonmessage.New()
	jmsg.SetInfo()
	jmsg.Message(msg...)
	return jmsg
}

// JSONWarnln is a shortcut function that will give you a JSON Message that is populated with time, level
// and your message ready to be shipped. It can still be decorated with more keys if needed.
func JSONWarnln(msg ...interface{}) *jsonmessage.JSONMessage {
	jmsg := jsonmessage.New()
	jmsg.SetWarn()
	jmsg.Message(msg...)
	return jmsg
}

// JSONCritln is a shortcut function that will give you a JSON Message that is populated with time, level
// and your message ready to be shipped. It can still be decorated with more keys if needed.
func JSONCritln(msg ...interface{}) *jsonmessage.JSONMessage {
	jmsg := jsonmessage.New()
	jmsg.SetCrit()
	jmsg.Message(msg...)
	return jmsg
}

// JSONDebugln is a shortcut function that will give you a JSON Message that is populated with time, level
// and your message ready to be shipped. It can still be decorated with more keys if needed.
func JSONDebugln(msg ...interface{}) *jsonmessage.JSONMessage {
	jmsg := jsonmessage.New()
	jmsg.SetDebug()
	jmsg.Message(msg...)
	return jmsg
}

// JSONInfof is a shortcut function that will give you a JSON Message that is populated with time, level
// and your message ready to be shipped. It can still be decorated with more keys if needed.
func JSONInfof(format string, vars ...interface{}) *jsonmessage.JSONMessage {
	jmsg := jsonmessage.New()
	jmsg.SetInfo()
	jmsg.Messagef(format, vars...)
	return jmsg
}

// JSONWarnf is a shortcut function that will give you a JSON Message that is populated with time, level
// and your message ready to be shipped. It can still be decorated with more keys if needed.
func JSONWarnf(format string, vars ...interface{}) *jsonmessage.JSONMessage {
	jmsg := jsonmessage.New()
	jmsg.SetWarn()
	jmsg.Messagef(format, vars...)
	return jmsg
}

// JSONCritf is a shortcut function that will give you a JSON Message that is populated with time, level
// and your message ready to be shipped. It can still be decorated with more keys if needed.
func JSONCritf(format string, vars ...interface{}) *jsonmessage.JSONMessage {
	jmsg := jsonmessage.New()
	jmsg.SetCrit()
	jmsg.Messagef(format, vars...)
	return jmsg
}

// JSONDebugf is a shortcut function that will give you a JSON Message that is populated with time, level
// and your message ready to be shipped. It can still be decorated with more keys if needed.
func JSONDebugf(format string, vars ...interface{}) *jsonmessage.JSONMessage {
	jmsg := jsonmessage.New()
	jmsg.SetDebug()
	jmsg.Messagef(format, vars...)
	return jmsg
}
