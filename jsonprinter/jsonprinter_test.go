package jsonprinter

import (
	"encoding/json"
	"regexp"
	"testing"

	"github.com/silverstagtech/gotracer"
	"github.com/silverstagtech/loggos/jsonmessage"
)

func TestJSONPrinter(t *testing.T) {
	tracing := gotracer.New()

	tests := []struct {
		name    string
		message string
	}{
		{
			name:    "first test",
			message: "one two three",
		},
		{
			name:    "second test",
			message: `$%^&*{}abcxyzABCXYZ@~/\`,
		},
	}

	for _, test := range tests {
		tracing.Reset()

		jp := New(10)
		jp.OverridePrinter(tracing)

		jmsg := jsonmessage.New()
		jmsg.Message(test.message)

		jp.Send(jmsg)
		<-jp.Flush()

		for _, jm := range tracing.Show() {
			jstruct := make(map[string]interface{})
			err := json.Unmarshal([]byte(jm), &jstruct)
			if err != nil {
				t.Logf("Failed to decode the json message. Error: %s", err)
				t.FailNow()
			}

			if test.message != jstruct[jsonmessage.JSONMessageKey] {
				t.Logf("Message in the tracer is not the message expected. Got: `%s`, Want: `%s`", test.message, jstruct[jsonmessage.JSONMessageKey])
				t.Fail()
			}
		}
	}
}

func TestPrettyPrinter(t *testing.T) {
	tracing := gotracer.New()

	testregexp := `^{\n\s{4}"`
	re := regexp.MustCompile(testregexp)

	jp := New(10)
	jp.OverridePrinter(tracing)
	jp.EnablePrettyPrint(true)

	jm := jsonmessage.New()
	jm.Message("Just a test")
	jm.SetInfo()

	jp.Send(jm)

	<-jp.Flush()

	if !re.MatchString(tracing.Show()[0]) {
		t.Log("Pretty print is not giving the correct signature.")
		t.Fail()
	}
}

func TestDebugPrinter(t *testing.T) {
	tracing := gotracer.New()

	tests := []struct {
		name          string
		turnOnDebug   bool
		expectMessage bool
	}{
		{
			name:          "On",
			turnOnDebug:   true,
			expectMessage: true,
		}, {
			name:          "off",
			turnOnDebug:   false,
			expectMessage: false,
		},
	}

	for _, test := range tests {
		jp := New(20)
		jp.OverridePrinter(tracing)

		tracing.Reset()
		jm := jsonmessage.New()
		jm.SetDebug()

		msgText := "Debug is off"

		if test.turnOnDebug {
			msgText = "Debug is on"
			jp.EnableDebugLogging(true)
		}

		jm.Message(msgText)
		jp.Send(jm)
		<-jp.Flush()

		if test.expectMessage {
			if tracing.Len() < 1 {
				t.Logf("Expected a debug message but didn't get one. Got: %s", tracing.Show())
				t.Fail()
			}
		} else {
			if tracing.Len() != 0 {
				t.Logf("Was not expecting a message but got one. Got: %s", tracing.Show())
				t.Fail()
			}
		}
	}
}

func TestShutdownSending(t *testing.T) {
	tracing := gotracer.New()

	jm := jsonmessage.New()
	jm.SetInfo()
	jm.Message("Test shutdown message.")

	jp := New(10)
	jp.OverridePrinter(tracing)
	jp.Send(jm)
	<-jp.Flush()

	if tracing.Len() != 1 {
		t.Log("Expecting a message but didn't get one.")
		t.Fail()
	}

	tracing.Reset()

	jp.Send(jm)
	<-jp.Flush()

	if tracing.Len() > 0 {
		t.Log("Expecting no messages but got something.")
		t.Fail()
	}
}
