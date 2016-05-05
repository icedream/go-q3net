package quake

type ReceivedMessage struct {
	*Message
	writer *Writer
}

func (m *ReceivedMessage) Writer() *Writer { return m.writer }
