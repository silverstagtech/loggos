package lineprinter

import (
	"fmt"
	"time"

	"github.com/silverstagtech/loggos/overrides"
	"github.com/silverstagtech/loggos/shared"
)

var (
	// DefaultLineTimeStampFunc is a time format override function for the line logger. Messages will have these prepended as soon as they arrive.
	DefaultLineTimeStampFunc func() string
)

func init() {
	DefaultLineTimeStampFunc = func() string { return time.Now().Format("Mon Jan _2 2006 15:04:05") }
}

// StandardLineLogger exposes the functions that match up with log interest levels.
type StandardLineLogger interface {
	Infoln(...interface{})
	Infof(string, ...interface{})
	Warnln(...interface{})
	Warnf(string, ...interface{})
	Critln(...interface{})
	Critf(string, ...interface{})
	Flush() chan bool
	OverrideTimeStamping(func() string)
	OverridePrinter(overrides.Overrider)
	EnableAuditMode(bool)
}

// DebugLineLogger uses StandardLogger but also includes Debugging logs.
type DebugLineLogger interface {
	StandardLineLogger
	Debugln(...interface{})
	Debugf(string, ...interface{})
	EnableDebugLogging(bool)
}

// Logger collects logs and prints them to the console in the order that it gets them.
// It needs to be flushed when the user if finished with to to not loose any logs.
type Logger struct {
	printDebug        bool
	logsToPrint       chan string
	FinishedChan      chan bool
	shutdown          bool
	transportOverride overrides.Overrider
	auditmode         bool
	droppedMessages   int64
	timestampFunc     func() string
}

// New creates a logger and returns it.
func New(buffer uint) *Logger {
	l := &Logger{
		logsToPrint:   make(chan string, buffer),
		FinishedChan:  make(chan bool, 1),
		timestampFunc: DefaultLineTimeStampFunc,
	}
	go l.printlogs()
	return l
}

// OverrideTimeStamping is used to change the default timestamp on individual loggers.
func (l *Logger) OverrideTimeStamping(f func() string) {
	l.timestampFunc = f
}

// OverridePrinter is used to insert your own function for hijacking the message on the
// way to the console. This allows you to push the log message to where ever you want.
func (l *Logger) OverridePrinter(override overrides.Overrider) {
	l.transportOverride = override
}

// EnableDebugLogging signals the Logger to print debug messages.
func (l *Logger) EnableDebugLogging(toggle bool) {
	l.printDebug = toggle
}

// EnableAuditMode will cause the logger to slow down if it us unable to process logs fast enough.
// Consider using this with a high buffer count.
func (l *Logger) EnableAuditMode(toggle bool) {
	l.auditmode = toggle
}

func (l *Logger) printlogs() {
	for {
		select {
		case msg, ok := <-l.logsToPrint:
			if !ok {
				l.FinishedChan <- true
				return
			}
			if l.transportOverride != nil {
				l.transportOverride.Send(msg)
			} else {
				l.defaultPrinter(msg)
			}
		}
	}
}

func (l *Logger) defaultPrinter(msg string) {
	fmt.Println(msg)
}

// Flush stops the logger from consuming more messages.
// Flush returns a chan bool to tell you when all messages
// have been printed. The channel will be closed once all messages have been flushed.
func (l *Logger) Flush() chan bool {
	l.shutdown = true
	close(l.logsToPrint)
	return l.FinishedChan
}

func (l *Logger) prepender(tag, msg string) string {
	if len(msg) == 0 {
		return fmt.Sprintf("%s %s", l.timestampFunc(), tag)
	}
	return fmt.Sprintf("%s %s %s", l.timestampFunc(), tag, msg)
}

func (l *Logger) prependInfo(msg string) string {
	return l.prepender(shared.InformationMessage, msg)
}

func (l *Logger) prependWarn(msg string) string {
	return l.prepender(shared.WarningMessage, msg)
}

func (l *Logger) prependCrit(msg string) string {
	return l.prepender(shared.CriticalMessage, msg)
}

func (l *Logger) prependDebug(msg string) string {
	return l.prepender(shared.DebugMessage, msg)
}

// Infoln takes a string adds a new line to the end and sends it to be printed
func (l *Logger) Infoln(msg ...interface{}) {
	if l.shutdown {
		return
	}
	out := []interface{}{l.prependInfo("")}
	out = append(out, msg...)
	l.send(fmt.Sprintln(out...))
}

// Warnln takes a string adds a new line to the end and sends it to be printed
func (l *Logger) Warnln(msg ...interface{}) {
	if l.shutdown {
		return
	}
	out := []interface{}{l.prependWarn("")}
	out = append(out, msg...)
	l.send(fmt.Sprintln(out...))
}

// Critln takes a string adds a new line to the end and sends it to be printed
func (l *Logger) Critln(msg ...interface{}) {
	if l.shutdown {
		return
	}
	out := []interface{}{l.prependCrit("")}
	out = append(out, msg...)
	l.send(fmt.Sprintln(out...))
}

// Debugln takes a string adds a new line to the end and sends it to be printed
func (l *Logger) Debugln(msg ...interface{}) {
	if l.shutdown {
		return
	}
	if !l.printDebug {
		return
	}
	out := []interface{}{l.prependDebug("")}
	out = append(out, msg...)
	l.send(fmt.Sprintln(out...))
}

// Infof takes a format string and as many vars as needed, merges the format with vars
// then sends the message to be printed
func (l *Logger) Infof(format string, vars ...interface{}) {
	if l.shutdown {
		return
	}
	l.send(l.prependInfo(fmt.Sprintf(format, vars...)))
}

// Warnf takes a format string and as many vars as needed, merges the format with vars
// then sends the message to be printed
func (l *Logger) Warnf(format string, vars ...interface{}) {
	if l.shutdown {
		return
	}
	l.send(l.prependWarn(fmt.Sprintf(format, vars...)))
}

// Critf takes a format string and as many vars as needed, merges the format with vars
// then sends the message to be printed
func (l *Logger) Critf(format string, vars ...interface{}) {
	if l.shutdown {
		return
	}
	l.send(l.prependCrit(fmt.Sprintf(format, vars...)))
}

// Debugf takes a format string and as many vars as needed, merges the format with vars
// then sends the message to be printed
func (l *Logger) Debugf(format string, vars ...interface{}) {
	if l.shutdown {
		return
	}
	if !l.printDebug {
		return
	}
	l.send(l.prependDebug(fmt.Sprintf(format, vars...)))
}

// send will select the correct sending function for shipping logs.
func (l *Logger) send(msg string) {
	if l.auditmode {
		shared.AuditSender(msg, l.logsToPrint)
		return
	}

	shared.BestEffortSender(msg, l.logsToPrint, l.droppedMessage)
}

func (l *Logger) droppedMessage() {
	l.droppedMessages++
}
