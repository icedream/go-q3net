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
	msg := append([]byte{0xff, 0xff, 0xff, 0xff}, []byte("getinfo")...)
	msg = append(msg, 0x20)
	msg = append(msg, []byte("xxx")...)
	msg = append(msg, 0x20)
	msg = append(msg, []byte("\"yy\\\"y\"")...)
	msg = append(msg, 0x20)
	msg = append(msg, []byte("'yz\"y'")...)
	msg = append(msg, 0x0A)

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
