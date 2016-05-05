package quake

import (
	"bytes"
	"io"
	"strings"
)

type Message struct {
	Header         *Header
	Name           string
	FragmentOffset uint32
	FragmentLength uint32
	Data           []byte
}

func UnmarshalMessage(data []byte) (*Message, error) {
	// Is there even a full header?
	if len(data) < 4 {
		return nil, ErrCorruptedMessage
	}

	header, err := UnmarshalHeader(data[0:4])
	if err != nil {
		return nil, err
	}

	// netchan UDP packets are structured differently
	if !header.IsOOB() {
		/*
			Original commentary from https://github.com/id-Software/Quake-III-Arena/blob/master/code/qcommon/net_chan.c#L30-L33:

			packet header
			-------------
			4	outgoing sequence.  high bit will be set if this is a fragmented message
			[2	qport (only for client to server)]
			[2	fragment start byte]
			[2	fragment length. if < FRAGMENT_SIZE, this is the last fragment]

			if the sequence number is -1, the packet should be handled as an out-of-band
			message instead of as part of a netcon.

			All fragments will have the same sequence numbers.

			The qport field is a workaround for bad address translating routers that
			sometimes remap the client's source port on a packet during gameplay.

			If the base part of the net address matches and the qport matches, then the
			channel matches even if the IP port differs.  The IP port should be updated
			to the new value before sending out any replies.

		*/

		// TODO - if packet is from client to server then 16 bits client's port number here

		if header.Fragmented {
			// TODO

			// 16 bits fragment offset

			// 16 bits fragment length
		}
	}

	splitPos := bytes.IndexAny(data[4:], " \n\r\t\x00\\")
	if splitPos < 0 {
		splitPos = len(data) - 4
	}
	commandName := string(data[4 : 4+splitPos])

	//log.Printf("Got message with command name %q", commandName)

	extra := []byte(nil)
	if data[4+splitPos] == 92 {
		splitPos--
	}
	if len(data) > 4+splitPos {
		extra = data[4+splitPos+1:]
	}
	for len(extra) > 0 && extra[len(extra)-1] == 0 {
		extra = extra[0 : len(extra)-1]
	}

	return &Message{
		Name:   commandName,
		Data:   extra,
		Header: header,
	}, nil
}

func (m *Message) SetArguments(argv []string) {
	m.Data = make([]byte, 0)

	for _, arg := range argv {
		if strings.Contains(arg, " ") {
			arg = "\"" + strings.Replace(arg, "\"", "\\\"", -1) + "\""
		}

		if len(m.Data) > 0 {
			// separate with space
			m.Data = append(m.Data, 0x20)
		}

		// add argument in binary form
		m.Data = append(m.Data, []byte(arg)...)
	}

	//m.Data = append(m.Data, 0x0A)
}

func (m *Message) GetArguments() []string {
	// All matches, extract arguments
	buffer := m.Data
	argv := []string{}
	for len(buffer) > 0 {
		c := rune(buffer[0])

		switch c {
		case '"', '\'':
			// quoted string
			searchPos := 1
			for {
				pos := bytes.IndexByte(buffer[searchPos:], byte(c))
				if pos < 0 {
					// found quote start without a matching end
					return nil
				}

				pos += searchPos

				if rune(buffer[pos-1]) != '\\' {
					// found end of quoted string
					searchPos = pos
					break
				}

				// that's an escaped quote char, skip that
				searchPos = pos + 1
			}

			// append what's inside the quotes to the arguments list
			arg := string(buffer[1:searchPos])
			arg = strings.Replace(arg, "\\"+string(c), string(c), -1)
			argv = append(argv, arg)

			if searchPos+1 < len(buffer) {
				buffer = buffer[searchPos+1:]
			} else {
				buffer = []byte{}
			}
		case ' ', '\t', '\n', '\r':
			buffer = buffer[1:]
		default:
			// search for next space
			pos := bytes.IndexByte(buffer, 0x20)
			if pos < 0 {
				pos = bytes.IndexByte(buffer, 0x0A)
			}
			if pos < 0 {
				pos = bytes.IndexByte(buffer, 0x00)
			}
			if pos < 0 {
				pos = len(buffer)
			}
			arg := string(buffer[0:pos])
			argv = append(argv, arg)

			if pos+1 < len(buffer) {
				buffer = buffer[pos+1:]
			} else {
				buffer = []byte{}
			}
		}

	}

	return argv
}

func (m *Message) Marshal(w io.Writer) error {
	// Buffer
	buf := new(bytes.Buffer)

	// Header
	if _, err := buf.Write(m.Header.Marshal()); err != nil {
		return err
	}

	// Command name
	if _, err := buf.Write([]byte(m.Name)); err != nil {
		return err
	}

	// Data
	if m.Data != nil && len(m.Data) > 0 {
		if _, err := buf.Write([]byte{0x20}); err != nil {
			return err
		}
	}
	if _, err := buf.Write(m.Data); err != nil {
		return err
	}

	// Separator, IW code expects this to be there, however our
	// implementation doesn't require it.
	buf.Write([]byte{0x00})

	// And now write as one message
	w.Write(buf.Bytes())

	return nil
}
