package quake

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGenerateMessage(t *testing.T) {
	actualMsgInstance := &Message{
		Header: OOBHeader,
		Name:   "getinfo",
	}
	actualMsgInstance.SetArguments([]string{"x\" y\"z"})
	actualMsgBuf := new(bytes.Buffer)
	if err := actualMsgInstance.Marshal(actualMsgBuf); err != nil {
		t.Fatal(err)
		return
	}
	actualMsg := actualMsgBuf.Bytes()

	expectedMsg := append([]byte{0xff, 0xff, 0xff, 0xff}, []byte("getinfo")...)
	expectedMsg = append(expectedMsg, 0x20)
	expectedMsg = append(expectedMsg, []byte("\"x\\\" y\\\"z\"")...)
	expectedMsg = append(expectedMsg, 0x00) // delimiter

	assert.EqualValues(t, expectedMsg, actualMsg)
}

func TestParseMessage(t *testing.T) {
	testWithDelimiter(t, 0x0A)
}

func TestParseMessage_ZeroDelimiter(t *testing.T) {
	testWithDelimiter(t, 0x00)
}

func TestParseMessage_DoubleZeroDelimiter(t *testing.T) {
	testWithDelimiter(t, 0x00, 0x00)
}

func testWithDelimiter(t *testing.T, delimiter ...byte) {
	msg := append([]byte{0xff, 0xff, 0xff, 0xff}, []byte("getinfo")...)
	msg = append(msg, 0x20)
	msg = append(msg, []byte("xxx")...)
	msg = append(msg, 0x20)
	msg = append(msg, []byte("\"yy\\\"y\"")...)
	msg = append(msg, 0x20)
	msg = append(msg, []byte("'yz\"y'")...)
	msg = append(msg, delimiter...)
	parsedMsg, err := UnmarshalMessage(msg)
	if err != nil {
		t.Errorf("Parser threw an error: %v", err)
	}
	switch parsedMsg.Name {
	case "getinfo":
		args := parsedMsg.GetArguments()
		if !assert.Equal(t, 3, len(args)) {
			t.Fail()
			return
		}
		assert.Equal(t, "xxx", args[0])
		assert.Equal(t, "yy\"y", args[1])
		assert.Equal(t, "yz\"y", args[2])
	default:
		t.Errorf("Unexpected parsed command name: %v", parsedMsg.Name)
	}
}

func TestParseMessage_Crash_Unknown1(t *testing.T) {
	msg := []byte("\xbd\xbf\x00G\xef\x85\xff\xf0\xbd\xbf\x16")
	parsedMsg, err := UnmarshalMessage(msg)
	if err != nil {
		t.Errorf("Parser threw an error: %v", err)
	}
	assert.NotEmpty(t, parsedMsg, "parsedMsg should not be empty")
}

func TestParseMessage_Nothing(t *testing.T) {
	msg := []byte{}
	parsedMsg, err := UnmarshalMessage(msg)
	if err == nil {
		t.Error("Parser did not throw an error on invalid message!")
	}
	assert.Empty(t, parsedMsg, "parsedMsg should be empty on error")
}

func TestParseMessage_InvalidHeaderWithoutContent(t *testing.T) {
	msg := []byte{'0', '0', '0', '0'}
	parsedMsg, err := UnmarshalMessage(msg)
	if err != nil {
		t.Errorf("Parser threw an error: %v", err)
	}
	assert.Empty(t, parsedMsg.Name, "command name must be empty")
}

func TestParseMessage_WithoutContent(t *testing.T) {
	msg := []byte{0xff, 0xff, 0xff, 0xff}
	parsedMsg, err := UnmarshalMessage(msg)
	if err != nil {
		t.Errorf("Parser threw an error: %v", err)
	}
	assert.Empty(t, parsedMsg.Name, "command name must be empty")
}

func TestParseMessage_WithoutArguments(t *testing.T) {
	msg := append([]byte{0xff, 0xff, 0xff, 0xff}, []byte("getinfo")...)
	parsedMsg, err := UnmarshalMessage(msg)
	if err != nil {
		t.Errorf("Parser threw an error: %v", err)
	}
	switch parsedMsg.Name {
	case "getinfo":
		args := parsedMsg.GetArguments()
		if !assert.Equal(t, 0, len(args)) {
			t.Fail()
			return
		}
	default:
		t.Errorf("Unexpected parsed command name: %v", parsedMsg.Name)
	}
}
