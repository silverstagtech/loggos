package loggos

import (
	"testing"

	"github.com/silverstagtech/gotracer"
)

func callAllLineFunctions(checkPanic bool, t *testing.T) {
	if checkPanic {
		panicChecker := func(t *testing.T) {
			if err := recover(); err != nil {
				t.Logf("Line printer Paniced. Error: %s", err)
				t.Fail()
			}
		}
		defer panicChecker(t)
	}

	Infoln("Test Message - Infoln")
	Warnln("Test Message - Warnln")
	Critln("Test Message - Critln")
	Debugln("Test Message - Debugln")
	Infof("Test Message - %s", "Infof")
	Warnf("Test Message - %s", "Warnf")
	Critf("Test Message - %s", "Critf")
	Debugf("Test Message - %s", "Debugf")
}

func TestMessageCount(t *testing.T) {
	// Call all line message functions with Debug on
	// Expect 8 messages
	// Call all line message functions with Debug off
	// Expect 6 messages

	tracing := gotracer.New()

	tests := []struct {
		name               string
		debugToggle        bool
		expectMessageCount int
	}{
		{
			name:               "debug on, expect 8 messages",
			debugToggle:        true,
			expectMessageCount: 8,
		},
		{
			name:               "debug of, expect 6 messages",
			debugToggle:        false,
			expectMessageCount: 6,
		},
	}

	for _, test := range tests {
		tracing.Reset()

		LineLoggerEnableDebugLogging(test.debugToggle)
		DefaultLineLogger.OverridePrinter(tracing)
		callAllLineFunctions(false, t)
		shutdownCurrentLoggers()

		if tracing.Len() != test.expectMessageCount {
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

func TestLinSendPanic(t *testing.T) {
	LineLoggerEnableDebugLogging(true)
	callAllLineFunctions(true, t)
}
