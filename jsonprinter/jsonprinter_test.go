package jsonprinter

import (
	"encoding/json"
	"fmt"
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

func TestDecorations(t *testing.T) {
	d1 := map[string]interface{}{
		"testing_key_1": 1,
		"testing_key_2": "two",
		"testing_key_3": map[string]string{
			"three": "number 3",
		},
	}

	d2 := map[string]interface{}{
		"bool_value": true,
	}

	tracing := gotracer.New()
	jp := New(100)
	jp.OverridePrinter(tracing)
	jp.EnablePrettyPrint(true)

	jp.AddDecoration(d1)
	jp.AddDecoration(d2)

	jm := jsonmessage.New()
	jm.SetInfo()
	jm.Message("Test decoration message.")

	jp.Send(jm)
	<-jp.Flush()

	regexMatches := []string{
		`"testing_key_1": 1`,
		`"three": "number 3"`,
		`bool_value": true`,
	}

	jmsg := tracing.Show()[0]

	for _, matcher := range regexMatches {
		if !regexp.MustCompile(matcher).MatchString(jmsg) {
			t.Logf("Failed to match %s decoration. Raw String:\n%s", matcher, jmsg)
			t.Fail()
		}
	}
}

func TestHumanTimestampping(t *testing.T) {
	tracing := gotracer.New()
	jp := New(100)
	jp.OverridePrinter(tracing)
	jp.EnablePrettyPrint(true)
	jp.EnableHumanTimestamps(true)

	jm := jsonmessage.New()
	jm.SetInfo()
	jm.Message("Test decoration message.")

	jp.Send(jm)
	<-jp.Flush()

	if !regexp.MustCompile(
		fmt.Sprintf(`"%s": "*"`, jsonmessage.JSONTimeStampKeyHuman),
	).MatchString(
		tracing.Show()[0],
	) {
		t.Logf("Did not see a human readable timestamp. Raw String:\n%s", tracing.Show()[0])
		t.Fail()
	}
}
