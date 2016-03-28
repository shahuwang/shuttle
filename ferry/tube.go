package ferry

import (
	"bufio"
	// "fmt"
	"net"
	"time"
)

type ITube interface {
	Send(data []byte) error
	Recieve() error
	Dispatch(conn *net.TCPConn)
}

type Tube struct {
	fromConn *net.TCPConn
	tunnel   *Tunnel
}

func (self *Tube) Recieve() (err error) {
	bufsize := PACKAGE_SIZE * 2
	rd := bufio.NewReaderSize(self.fromConn, bufsize)
	defer self.fromConn.CloseRead()
	for {
		buffer := BufferPool()
		n, err := rd.Read(buffer)
		if err != nil {
			// TODO
			break
		}
		err = self.Send(buffer[:n])
		if err != nil {
			break
		}
	}
	return nil
}

func (self *Tube) Send(data []byte) (err error) {
	for {
		_, err := self.tunnel.Write(data)
		if err != nil {
			//TODO
			break
		}
	}
	return nil
}

func (self *Tube) Dispatch(conn *net.TCPConn) {
	conn.SetKeepAlive(true)
	conn.SetKeepAlivePeriod(time.Second * 60)
	self.fromConn = conn
	self.Recieve()
}
