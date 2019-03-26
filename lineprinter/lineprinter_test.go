package lineprinter

import (
	"regexp"
	"testing"

	"github.com/silverstagtech/gotracer"
)

func TestLoggerOverride(t *testing.T) {
	tracing := gotracer.New()

	logger := New(10)
	logger.OverridePrinter(tracing)
	logger.Infoln("test message")
	finshed := logger.Flush()
	<-finshed

	if tracing.Len() == 0 {
		t.Logf("Logger override did not work as expected.")
		t.Fail()
	}
}

func TestShutdownLogger(t *testing.T) {
	tracing := gotracer.New()

	logger := New(10)
	logger.EnableDebugLogging(true)
	logger.OverridePrinter(tracing)
	// Finish early, no more messages should be able to be printed now.
	<-logger.Flush()

	// Try to print messages.
	logger.Infoln("test message")
	logger.Warnln("testmessage")
	logger.Critln("test message")
	logger.Debugln("test message")
	logger.Infof("test message")
	logger.Warnf("test message")
	logger.Critf("test message")
	logger.Debugf("test message")

	if tracing.Len() != 0 {
		t.Logf("Logger accepted messages after being flushed. %v", tracing.Show())
		t.Fail()
	}
}

func TestPrintDebugWhenDisabled(t *testing.T) {
	tracing := gotracer.New()

	logger := New(10)
	logger.OverridePrinter(tracing)
	logger.Debugln("test message")
	logger.Debugf("test message")

	if tracing.Len() != 0 {
		t.Logf("Logger accepted debug messages before being enabled. %v", tracing.Show())
		t.Fail()
	}

	logger.EnableDebugLogging(true)
	logger.Debugln("test message")
	logger.Debugf("test message")
	logger.EnableDebugLogging(false)
	logger.EnableDebugLogging(true)

	if tracing.Len() > 2 {
		t.Logf("Logger accepted debug messages after being disabled. %v", tracing.Show())
		t.Fail()
	}
}

func TestLoggerAppends(t *testing.T) {
	tracing := gotracer.New()

	tests := []struct {
		testName string
		log      string
		match    string
		exec     string
	}{
		{
			testName: "Info Messages",
			log:      "test message",
			match:    `INFO test message(\r\n|\r|\n)$`,
			exec:     "i",
		},
		{
			testName: "Warning Messages",
			log:      "test message",
			match:    `WARN test message(\r\n|\r|\n)$`,
			exec:     "w",
		},
		{
			testName: "Critical Messages",
			log:      "test message",
			match:    `CRIT test message(\r\n|\r|\n)$`,
			exec:     "c",
		},
		{
			testName: "Debug Messages",
			log:      "test message",
			match:    `DEBUG test message(\r\n|\r|\n)$`,
			exec:     "d",
		},
	}

	for _, funcLabel := range []string{"f", "ln"} {
		for _, test := range tests {
			logger := New(10)
			logger.OverridePrinter(tracing)
			tracing.Reset()

			switch test.exec {
			case "i":
				if funcLabel == "f" {
					logger.Infof("%s\n", test.log)
				}
				if funcLabel == "ln" {
					logger.Infoln(test.log)
				}
			case "w":
				if funcLabel == "f" {
					logger.Warnf("%s\n", test.log)
				}
				if funcLabel == "ln" {
					logger.Warnln(test.log)
				}
			case "c":
				if funcLabel == "f" {
					logger.Critf("%s\n", test.log)
				}
				if funcLabel == "ln" {
					logger.Critln(test.log)
				}
			case "d":
				logger.EnableDebugLogging(true)
				if funcLabel == "f" {
					logger.Debugf("%s\n", test.log)
				}
				if funcLabel == "ln" {
					logger.Debugln(test.log)
				}
			}

			finshed := logger.Flush()
			<-finshed

			re := regexp.MustCompile(test.match)
			for _, lgm := range tracing.Show() {
				if !re.Match([]byte(lgm)) {
					t.Logf("Log message does not start start with the correct label. Using %s function.\nGot: %s\nShould Match: %s", funcLabel, lgm, test.match)
					t.Fail()
				}
			}
		}
	}
}

func TestTimeStampFunc(t *testing.T) {
	tracing := gotracer.New()

	tests := []struct {
		name          string
		testRegex     string
		datestampfunc func() string
	}{
		{
			name:          "Static date",
			testRegex:     "^--static-- ",
			datestampfunc: func() string { return "--static--" },
		},
		{
			name:          "Default date",
			testRegex:     `(?:[A-z][a-z][a-z]\s){2}\s?(?:[0-9]|[0-9][0-9])\s[0-9]{4}\s(?:[0-9]+:?){3}`,
			datestampfunc: nil,
		},
	}

	for _, test := range tests {
		logger := New(10)
		tracing.Reset()
		if test.datestampfunc != nil {
			logger.OverrideTimeStamping(test.datestampfunc)
		}
		logger.OverridePrinter(tracing)
		logger.Critf("test message")
		<-logger.Flush()

		re := regexp.MustCompile(test.testRegex)

		if !re.MatchString(tracing.Show()[0]) {
			t.Logf("Date stamp function did not inject the correct string.")
			t.Fail()
		}
	}
}
