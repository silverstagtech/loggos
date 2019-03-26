package jsonmessage

import (
	"fmt"
	"time"
)

// JSONTimeStamper is used to stamp a JSON message with a time stamp of the users desire
type JSONTimeStamper interface {
	Stamp() string
}

type timeStamper struct{}

func newStamper() *timeStamper {
	return &timeStamper{}
}

// timeStampEpochNano returns the time now to be used as a time stamp.
func (s *timeStamper) timeStampEpochNano() string {
	return fmt.Sprintf("%d", time.Now().UnixNano())
}

func (s *timeStamper) Stamp() string {
	return s.timeStampEpochNano()
}
