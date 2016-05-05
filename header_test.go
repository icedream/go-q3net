package quake

import (
	"bytes"
	"strconv"
	"strings"
	"testing"
)

func TestNetchanConstants(t *testing.T) {
	binaryConst := strings.TrimLeft(strconv.FormatUint(uint64(FRAGMENT_BIT), 2), "0")
	if binaryConst != "1"+strings.Repeat("0", 31) {
		t.Fatalf("FRAGMENT_BIT is incorrect (value currently is %v as binary)", binaryConst)
	}

	binaryConst = strings.TrimLeft(strconv.FormatUint(uint64(MAX_MSGSEQNUM), 2), "0")
	if binaryConst != strings.Repeat("1", 31) {
		t.Fatalf("MAX_MSGSEQNUM is incorrect (value currently is %v as binary)", binaryConst)
	}
}

func TestHeaderIsOOB(t *testing.T) {
	h := &Header{
		Fragmented:            true,
		MessageSequenceNumber: MAX_MSGSEQNUM,
	}
	if !h.IsOOB() {
		t.Fatal("Header unexpectedly does not report itself as out of band header")
	}
}

func TestHeaderIsOOBConst(t *testing.T) {
	h := OOBHeader
	if !h.IsOOB() {
		t.Fatal("Header unexpectedly does not report itself as out of band header")
	}
}

func TestHeaderUnmarshalOOB(t *testing.T) {
	h, err := UnmarshalHeader([]byte{0xFF, 0xFF, 0xFF, 0xFF})
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	if !h.IsOOB() {
		t.Fatal("Header unexpectedly does not report itself as out of band header")
	}
}

func TestHeaderUnmarshal(t *testing.T) {
	h, err := UnmarshalHeader([]byte{0x01, 0x00, 0x00, 0x00})
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	if h.Fragmented {
		t.Fatal("Header unexpectedly does report fragmented")
	}
	if h.MessageSequenceNumber != 1 {
		t.Fatal("Message sequence number mismatch")
	}
}

func TestHeaderMarshal(t *testing.T) {
	h := &Header{
		Fragmented:            false,
		MessageSequenceNumber: 1,
	}
	data := h.Marshal()
	if !bytes.Equal(data, []byte{0x01, 0x00, 0x00, 0x00}) {
		t.Fatalf("Marshaled data is invalid, result is %x", data)
	}
}
