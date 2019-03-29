package jsonmessage

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/silverstagtech/loggos/shared"
)

var (
	// JSONTimeStampKey is used to write the time stamp in json messages
	JSONTimeStampKey = "timestamp"
	// JSONLevelKey is used to write the log message concern level.
	JSONLevelKey = "level"
	// JSONMessageKey is used as the log message key in the JSON structure
	JSONMessageKey = "log_message"
	// JSONTimeStampFunc is a override function to set the timestamp to what the user wants.
	// It must retrun a string which allows many formats to fit in.
	// By default the time stamp is a epoch nano number
	// Can be set by a init function in a higher level package.
	JSONTimeStampFunc JSONTimeStamper
	// JSONTimeStampKeyHuman is used when writing a human readable timestamp to the message
	JSONTimeStampKeyHuman = "human_readable_timestamp"
	// HumanTimeStampFormat is the time stamp format used for the default timestamp formating
	HumanTimeStampFormat = "Mon Jan _2 2006 15:04:05 MST"
	// JSONErrorKey is used as the key for adding an error message
	JSONErrorKey = "error"
)

// JSONMessage is a structure that will contain the message that you want to send.
type JSONMessage struct {
	msg map[string]interface{}
}

// New returns a empty JSONMessage ready to be populated. The timestamp will have already been
// written.
func New() *JSONMessage {
	jm := &JSONMessage{
		msg: make(map[string]interface{}),
	}
	jm.msg[JSONTimeStampKey] = timeStamp()
	return jm
}

func timeStamp() string {
	if JSONTimeStampFunc == nil {
		JSONTimeStampFunc = newStamper()
	}

	return JSONTimeStampFunc.Stamp()
}

// Add adds on the key that you want to add your message. This can be anything you want.
func (j *JSONMessage) Add(key string, value interface{}) {
	j.msg[key] = value
}

// Addf acts like fmt.Sprintf and stores the result under the supplied key. The stored result will be a string.
func (j *JSONMessage) Addf(key string, format string, values ...interface{}) {
	j.msg[key] = fmt.Sprintf(format, values...)
}

// Errorf is a shortcut function that is equivalent to calling Addf(JSONErrorKey, format, values...).
func (j *JSONMessage) Errorf(format string, values ...interface{}) {
	j.Addf(JSONErrorKey, format, values...)
}

// Error is a shortcut function that is equivalent to calling Add(JSONErrorKey, err.Error()).
func (j *JSONMessage) Error(err error) {
	j.Add(JSONErrorKey, err.Error())
}

// SetInfo sets level to INFO
func (j *JSONMessage) SetInfo() {
	j.Add(JSONLevelKey, shared.InformationMessage)
}

// SetWarn sets level to WARN
func (j *JSONMessage) SetWarn() {
	j.Add(JSONLevelKey, shared.WarningMessage)
}

// SetCrit sets level to CRIT
func (j *JSONMessage) SetCrit() {
	j.Add(JSONLevelKey, shared.CriticalMessage)
}

// SetDebug sets level to DEBUG
func (j *JSONMessage) SetDebug() {
	j.Add(JSONLevelKey, shared.DebugMessage)
}

// Message sets the message key with what you pass in. It works like fmt.Sprint.
func (j JSONMessage) Message(m ...interface{}) {
	j.addMessage(fmt.Sprint(m...))
}

// Messagef sets the message key with what you pass in but also allows for string formatting.
// It works like fmt.Sprintf.
func (j JSONMessage) Messagef(format string, m ...interface{}) {
	j.addMessage(fmt.Sprintf(format, m...))
}

func (j *JSONMessage) addMessage(msg string) {
	j.Add(JSONMessageKey, msg)
}

// Bytes returns the []byte representation of your message. If there is a error decoding your message
// it will still return a []byte but will contain the error message in the form of {"Error": "message"}.
func (j *JSONMessage) Bytes() []byte {
	var b []byte
	b, err := json.Marshal(j.msg)
	if err != nil {
		b = []byte(fmt.Sprintf(`{"error": "%s","raw_string":"%v"}`, err, j.msg))
	}
	return b
}

// String returns the string presentation of your message. If there is a error decoding your message
// it will still return a string but will contain the error message in the form of {"Error": "message"}.
func (j *JSONMessage) String() string {
	return string(j.Bytes())
}

// PrettyBytes returns the []byte representation of your message with some json indentation to make it east to read.
// If there is a error decoding your message it will still return a []byte but will contain the error message
// in the form of {\n    "Error": "message"\n}.
func (j *JSONMessage) PrettyBytes() []byte {
	var b []byte
	b, err := json.MarshalIndent(j.msg, "", "    ")
	if err != nil {
		b = []byte(fmt.Sprintf(`{\n    "error": "%s",\n"raw_string": "%v"\n}`, err, j.msg))
	}
	return b
}

// PrettyString returns the string presentation of your message. If there is a error decoding your message
// it will still return a string but will contain the error message in the form of {\n    "Error": "message"\n}.
func (j *JSONMessage) PrettyString() string {
	return string(j.PrettyBytes())
}

// IsDebug will return a bool which will indicate that the message is a debug message.
func (j *JSONMessage) IsDebug() bool {
	if v, ok := j.msg[JSONLevelKey].(string); ok {
		if v == shared.DebugMessage {
			return true
		}
	}
	return false
}

// AddHumanTimestamp is used to convert epoch with nanoseconds to a human timestamp.
// The result is written to the message using the key save in JSONTimeStampKeyHuman.
// The key under JSONTimeStampKey MUST be a int64 or the key is simply not written.
// The key under JSONTimeStampKey MUST be a epoch with nano seconds timestamp eg. time.Now().UnixNano()
// else you will get a strange looking date.
func (j *JSONMessage) AddHumanTimestamp() {
	// Can we parse to a string?
	t, ok := j.msg[JSONTimeStampKey].(string)
	if !ok {
		return
	}
	// Is the string a number?
	epoch, err := strconv.ParseInt(t, 10, 64)
	if err != nil {
		return
	}

	// The number must be epoch with nano eg: time.Now().UnixNano()
	// else you will get a strange looking date.
	m := epoch / (int64(time.Second) / int64(time.Nanosecond))
	n := epoch % (int64(time.Second) / int64(time.Nanosecond))
	j.Add(JSONTimeStampKeyHuman, time.Unix(m, n).Format(HumanTimeStampFormat))
}

// RawDump return a pointer to a JSONMessages internal data structure.
// Used in conjunction with JSONPrinter Mutator.
func (j *JSONMessage) RawDump() map[string]interface{} {
	return j.msg
}
