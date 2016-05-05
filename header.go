package quake

import (
	"bytes"
	"encoding/binary"
)

const (
	FRAGMENT_BIT  = uint32(1) << 31
	MAX_MSGSEQNUM = ^FRAGMENT_BIT
)

var (
	endianness = binary.LittleEndian

	OOBHeader = &Header{
		Fragmented:            true,
		MessageSequenceNumber: MAX_MSGSEQNUM,
	}
)

type Header struct {
	Fragmented            bool
	MessageSequenceNumber uint32
}

func (h *Header) IsOOB() bool {
	return h.Fragmented && h.MessageSequenceNumber == MAX_MSGSEQNUM
}

func (h *Header) Marshal() []byte {
	// Fragmented is first bit
	// Message sequence number is the last 31 bits
	finalHeader := h.MessageSequenceNumber
	if h.Fragmented {
		finalHeader |= FRAGMENT_BIT
	}

	buf := new(bytes.Buffer)
	binary.Write(buf, endianness, finalHeader)
	bytes := buf.Bytes()

	return bytes
}

func UnmarshalHeader(data []byte) (*Header, error) {
	if len(data) < 4 {
		return nil, ErrCorruptedMessage
	}

	// Sequence number
	seqNum := endianness.Uint32(data)

	// First bit = Fragmented?
	fragmented := seqNum&FRAGMENT_BIT != 0
	if fragmented {
		seqNum &= ^FRAGMENT_BIT
	}

	return &Header{
		Fragmented:            fragmented,
		MessageSequenceNumber: seqNum,
	}, nil
}
