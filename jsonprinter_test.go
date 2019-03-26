package loggos

import (
	"testing"

	"github.com/silverstagtech/gotracer"
	"github.com/silverstagtech/loggos/jsonmessage"
)

func sendOnAllJSONFunctions(checkPanic bool, t *testing.T) {
	if checkPanic {
		panicChecker := func(t *testing.T) {
			if err := recover(); err != nil {
				t.Logf("JSONPrinter Paniced. Error: %s", err)
				t.Fail()
			}
		}
		defer panicChecker(t)
	}

	jms := []*jsonmessage.JSONMessage{}
	jms = append(jms, JSONInfoln("Test Message - JSONInfoln"))
	jms = append(jms, JSONWarnln("Test Message - JSONWarnln"))
	jms = append(jms, JSONCritln("Test Message - JSONCritln"))
	jms = append(jms, JSONDebugln("Test Message - JSONDebugln"))
	jms = append(jms, JSONInfof("Test Message - %s", "JSONInfof"))
	jms = append(jms, JSONWarnf("Test Message - %s", "JSONWarnf"))
	jms = append(jms, JSONCritf("Test Message - %s", "JSONCritf"))
	jms = append(jms, JSONDebugf("Test Message - %s", "JSONDebugf"))

	for _, m := range jms {
		SendJSON(m)
	}
}

func TestMessageCountJSON(t *testing.T) {
	tracing := gotracer.New()

	tests := []struct {
		name               string
		debugToggle        bool
		expectMessageCount int
	}{
		{
			name:               "Debug on, expect 8 messages",
			debugToggle:        true,
			expectMessageCount: 8,
		},
		{
			name:               "Debug off, expect 6 messages",
			debugToggle:        false,
			expectMessageCount: 6,
		},
	}

	for _, test := range tests {
		tracing.Reset()
		JSONLoggerEnableDebugLogging(test.debugToggle)
		DefaultJSONLogger.OverridePrinter(tracing)
		sendOnAllJSONFunctions(false, t)
		shutdownCurrentLoggers()

		if test.expectMessageCount != tracing.Len() {
			t.Logf(
				"%s - wanted %d messages but got %d. Messages:\n%s",
				test.name,
				test.expectMessageCount,
				tracing.Len(),
				tracing.Show(),
			)
			t.Fail()
		}
	}
}

func TestJSONSendPanic(t *testing.T) {

	JSONLoggerEnableDebugLogging(true)
	JSONLoggerEnablePrettyPrint(true)
	sendOnAllJSONFunctions(true, t)
	shutdownCurrentLoggers()
}
