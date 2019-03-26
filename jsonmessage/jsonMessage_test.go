package jsonmessage

import (
	"regexp"
	"testing"

	"github.com/silverstagtech/loggos/shared"
)

func TestLevelSettings(t *testing.T) {
	jm := New()

	jm.SetInfo()
	if jm.msg[JSONLevelKey] != shared.InformationMessage {
		t.Logf("SetInfo did not set the correct level, Got: %s", jm.msg[JSONLevelKey])
		t.Fail()
	}

	jm.SetWarn()
	if jm.msg[JSONLevelKey] != shared.WarningMessage {
		t.Logf("SetWarn did not set the correct level, Got: %s", jm.msg[JSONLevelKey])
		t.Fail()
	}

	jm.SetCrit()
	if jm.msg[JSONLevelKey] != shared.CriticalMessage {
		t.Logf("SetCrit did not set the correct level, Got: %s", jm.msg[JSONLevelKey])
		t.Fail()
	}

	jm.SetDebug()
	if jm.msg[JSONLevelKey] != shared.DebugMessage {
		t.Logf("SetDebug did not set the correct level, Got: %s", jm.msg[JSONLevelKey])
		t.Fail()
	}
}

func TestTimeStamping(t *testing.T) {
	jm := New()
	if _, ok := jm.msg[JSONTimeStampKey]; !ok {
		t.Logf("A new JSON message has no timestamp. Raw data: %v", jm.msg)
		t.Fail()
	}
}

type testTimeStamp struct {
	testString string
}

func (ts *testTimeStamp) Stamp() string { return ts.testString }

func TestCustomTimeStamp(t *testing.T) {
	testString := "gofer"
	JSONTimeStampFunc = &testTimeStamp{testString: testString}
	jm := New()

	if jm.msg[JSONTimeStampKey] != testString {
		t.Logf("Setting a custom timestamp function do not yet correct result.")
		t.Fail()
	}
	JSONTimeStampFunc = nil
}

func TestAddingCustomField(t *testing.T) {
	jm := New()

	key := "test"
	value := "Test Message"
	jm.Add(key, value)

	if jm.msg[key] != value {
		t.Logf("Adding a custom field to the json message didn't get added to the message.")
		t.FailNow()
	}

	key1 := "k1"
	key2 := "k2"
	nestedValue := "value"
	jm.Add(key1, map[string]string{key2: nestedValue})

	// This is a bit hairy but we need to cast to a map[string]string
	// as we get it as a map[string]interface{}
	if jm.msg[key1].(map[string]string)[key2] != nestedValue {
		t.Logf("Adding a nested string did not work. Raw data: %v", jm)
	}
}

func TestMessage(t *testing.T) {
	testMessage := "This is a test message"
	jm := New()
	jm.Message(testMessage)

	if jm.msg[JSONMessageKey] != testMessage {
		t.Log("Message failed to write the log message to the JSON structure.")
		t.Fail()
	}

	jm.Messagef("%s", testMessage)
	if jm.msg[JSONMessageKey] != testMessage {
		t.Log("Messagef failed to write the log message to the JSON structure.")
		t.Fail()
	}
}

func TestString(t *testing.T) {
	testMessage := "This is a test message"
	jm := New()
	jm.Message(testMessage)
	re := regexp.MustCompile(`"log_message":"This is a test message"`)
	if ok := re.MatchString(jm.String()); !ok {
		t.Log("Write the JSON message to a string did not return the correct message")
		t.Fail()
	}
}

func TestPrettyString(t *testing.T) {
	testMessage := "This is a test message"
	jm := New()
	jm.Message(testMessage)
	re := regexp.MustCompile(`\n\s+"log_message": "This is a test message",\n`)
	if ok := re.MatchString(jm.PrettyString()); !ok {
		t.Logf("PrettyString gave back the wrong message. Raw data: %v", jm.PrettyString())
		t.Fail()
	}
}
