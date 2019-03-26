package jsonprinter

import (
	"fmt"

	"github.com/silverstagtech/loggos/jsonmessage"
	"github.com/silverstagtech/loggos/overrides"
	"github.com/silverstagtech/loggos/shared"
)

// JSONLogger is a logger that implements the functions of this package
type JSONLogger interface {
	EnablePrettyPrint(bool)
	EnableAuditMode(bool)
	OverridePrinter(overrides.Overrider)
	Send(*jsonmessage.JSONMessage)
	Flush() chan bool
}

// DebugJSONLogger allowed you to also toggle debug messages on and off while also pulling in JSONLogger
type DebugJSONLogger interface {
	JSONLogger
	EnableDebugLogging(bool)
}

// JSONPrinter consumes JSON Logs and sends them to the current output
type JSONPrinter struct {
	printDebug        bool
	logsToPrint       chan string
	FinishedChan      chan bool
	shutdown          bool
	printPretty       bool
	auditmode         bool
	droppedMessages   int64
	transportOverride overrides.Overrider
}

// New created a empty JSON Printer and starts the printer ready for messages.
func New(buffer uint) *JSONPrinter {
	jp := &JSONPrinter{
		logsToPrint:  make(chan string, buffer),
		FinishedChan: make(chan bool, 1),
	}
	go jp.printlogs()
	return jp
}

// EnableDebugLogging signals the Logger to print debug messages.
func (j *JSONPrinter) EnableDebugLogging(toggle bool) {
	j.printDebug = toggle
}

// EnablePrettyPrint signals the Logger to print human readable messages.
func (j *JSONPrinter) EnablePrettyPrint(toggle bool) {
	j.printPretty = toggle
}

// EnableAuditMode will cause the logger to slow down if it us unable to process logs fast enough.
// Consider using this with a high buffer count.
func (j *JSONPrinter) EnableAuditMode(toggle bool) {
	j.auditmode = toggle
}

func (j *JSONPrinter) printlogs() {
	for {
		select {
		case msg, ok := <-j.logsToPrint:
			if !ok {
				j.FinishedChan <- true
				return
			}
			if j.transportOverride != nil {
				j.transportOverride.Send(msg)
			} else {
				j.defaultPrinter(msg)
			}
		}
	}
}

func (j JSONPrinter) defaultPrinter(msg string) {
	fmt.Println(msg)
}

// OverridePrinter is used to insert your own function for hijacking the message on the
// way to the console. This allows you to push the log message to where ever you want.
func (j *JSONPrinter) OverridePrinter(override overrides.Overrider) {
	j.transportOverride = override
}

// Flush stops the logger from consuming more messages.
// Flush returns a chan bool to tell you when all messages
// have been printed. The channel will be closed once all messages have been flushed.
func (j *JSONPrinter) Flush() chan bool {
	if j.shutdown {
		c := make(chan bool, 1)
		close(c)
		return c
	}

	close(j.logsToPrint)
	j.shutdown = true
	return j.FinishedChan
}

// Send takes a pointer to a JSONMessage and send it to the printer.
// If the logger is already shutdown then it will just silently consume the message.
func (j *JSONPrinter) Send(msg *jsonmessage.JSONMessage) {
	if j.shutdown {
		return
	}
	if msg.IsDebug() {
		if !j.printDebug {
			return
		}
	}

	if j.printPretty {
		j.send(msg.PrettyString())
		return
	}

	j.send(msg.String())
}

// send will select the correct sending function for shipping logs.
func (j *JSONPrinter) send(msg string) {
	if j.auditmode {
		shared.AuditSender(msg, j.logsToPrint)
		return
	}

	shared.BestEffortSender(msg, j.logsToPrint, j.droppedMessage)
}

func (j *JSONPrinter) droppedMessage() {
	j.droppedMessages++
}
