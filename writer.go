package quake

import (
	"bytes"
	"net"
)

type Writer struct {
	socket *net.UDPConn
	addr   *net.UDPAddr
}

/*func (w *Writer) Socket() *net.UDPConn {
	return w.socket
}*/

func (w *Writer) Addr() *net.UDPAddr {
	return w.addr
}

func (w *Writer) Write(msg *Message) error {
	buf := new(bytes.Buffer)
	msg.Marshal(buf)
	_, err := w.socket.WriteToUDP(buf.Bytes(), w.addr)
	return err
}
