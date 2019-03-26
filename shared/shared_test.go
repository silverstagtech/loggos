package shared

import (
	"testing"
	"time"
)

func TestBestEffortSender(t *testing.T) {
	c := make(chan string, 1)
	dropped := 0
	droppedFunc := func() {
		dropped++
	}

	// Send messages
	BestEffortSender("1", c, droppedFunc)
	// It should be full now.
	BestEffortSender("2", c, droppedFunc)

	if dropped == 0 {
		t.Logf("BestEffortSender did not drop messages when the queue is full.")
		t.Fail()
	}
}

func TestAuditSender(t *testing.T) {
	c := make(chan string, 1)

	// Send first message. It should accept this one.
	AuditSender("1", c)
	// Now its is full, so we need to send and timeout
	auditSenderWaiter := func() chan bool {
		senderChan := make(chan bool, 1)
		go func() {
			AuditSender("2", c)
			senderChan <- true
		}()
		return senderChan
	}
	ticker := time.NewTicker(time.Millisecond * 2)

	select {
	case <-auditSenderWaiter():
		t.Logf("AuditSender return when the channel was full.")
		t.FailNow()
	case <-ticker.C:
	}
}
