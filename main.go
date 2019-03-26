package loggos

import (
	"sync"

	"github.com/silverstagtech/loggos/jsonprinter"
	"github.com/silverstagtech/loggos/lineprinter"
)

var (
	// DefaultLineLogger is a logger that you can just call upon when you start your application.
	DefaultLineLogger lineprinter.DebugLineLogger
	// DefaultJSONLogger is a logger that you can just call upon when you start your application.
	DefaultJSONLogger jsonprinter.DebugJSONLogger
	// DefaultLineLoggerBuffer holds how many lines the line buffer will hold onto before it starts to drop
	// or slow down the application.
	DefaultLineLoggerBuffer = 500
	// DefaultJSONLoggerBuffer holds how many json messages the line buffer will hold onto before it starts to drop
	// or slow down the application.
	DefaultJSONLoggerBuffer = 500
)

func startdefaultLineLogger() {
	if DefaultLineLogger == nil {
		DefaultLineLogger = lineprinter.New(uint(DefaultLineLoggerBuffer))
	}
}

func startdefaultJSONLogger() {
	if DefaultJSONLogger == nil {
		DefaultJSONLogger = jsonprinter.New(uint(DefaultJSONLoggerBuffer))
	}
}

// Flush on the package will stop the default logging engines that you have made use of.
// It will close the channel once all logging engines have flushed everything and stopped.
func Flush() chan bool {
	c := make(chan bool, 1)

	go func() {
		wg := &sync.WaitGroup{}

		if DefaultJSONLogger != nil {
			wg.Add(1)
			go func() {
				<-DefaultJSONLogger.Flush()
				wg.Done()
			}()
		}

		if DefaultLineLogger != nil {
			wg.Add(1)
			go func() {
				<-DefaultLineLogger.Flush()
				wg.Done()
			}()
		}

		wg.Wait()
		close(c)
	}()

	return c
}
