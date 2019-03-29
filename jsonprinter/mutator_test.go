package jsonprinter

import (
	"regexp"
	"testing"

	"github.com/silverstagtech/gotracer"
	"github.com/silverstagtech/loggos/jsonmessage"
)

type TestMutator struct {
	mutator func(*jsonmessage.JSONMessage) bool
}

func (tm *TestMutator) Mutate(jm *jsonmessage.JSONMessage) bool {
	return tm.mutator(jm)
}

func TestGoodMutation(t *testing.T) {
	tracing := gotracer.New()
	jp := New(10)
	jp.OverridePrinter(tracing)

	jp.AddMutator(
		&TestMutator{
			mutator: func(jm *jsonmessage.JSONMessage) bool {
				internalMessage := jm.RawDump()
				internalMessage["magic_here"] = "Magic Everywhere"
				internalMessage["timestamp"] = "overridden"
				delete(internalMessage, "level")
				return true
			},
		},
	)

	jm := jsonmessage.New()
	jm.Message("Just a test")
	jm.SetInfo()
	jp.Send(jm)
	<-jp.Flush()

	messages := tracing.Show()
	if len(messages) < 1 {
		t.Logf("Tracing has no messages in it which means the mutation failed")
		t.FailNow()
	}

	if !regexp.MustCompile(`magic_here`).MatchString(messages[0]) {
		t.Logf("Mutator did not add in the magic_here key. Raw Message:\n%s", messages[0])
		t.FailNow()
	}

	if regexp.MustCompile(`level`).MatchString(messages[0]) {
		t.Logf("Mutator did not remove level key. Raw Message:\n%s", messages[0])
		t.FailNow()
	}
}

func TestBadMutation(t *testing.T) {
	tracing := gotracer.New()
	jp := New(10)
	jp.OverridePrinter(tracing)
	jp.AddMutator(
		&TestMutator{
			mutator: func(jm *jsonmessage.JSONMessage) bool {
				return false
			},
		},
	)

	jm := jsonmessage.New()
	jm.Message("Just a test")
	jm.SetInfo()
	jp.Send(jm)
	<-jp.Flush()
	if tracing.Len() > 0 {
		t.Logf("A failed mutator still rendered a message. Messages available:\n%v", tracing.Show())
	}
}
