package ferry

import (
	"bufio"
	"fmt"
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
			fmt.Println(err)
			fmt.Println("Recieve error")
			break
		}
		go self.Send(buffer[:n])
	}
	return err
}

func (self *Tube) ServerRecieve() (err error) {
	for {
		buffer := BufferPool()
		n, err := self.tunnel.Read(buffer)
		if err != nil {
			// TODO
			fmt.Println(err)
			fmt.Println("Recieve error")
			break
		}
		go self.ServerSend(buffer[:n])
	}
	return err

}

func (self *Tube) ServerSend(data []byte) (err error) {
	_, err = self.fromConn.Write(data)
	if err != nil {
		fmt.Println("server send err")
		fmt.Println(err)
	}
	return
}

func (self *Tube) Send(data []byte) (err error) {
	_, err = self.tunnel.Send(data)
	if err != nil {
		//TODO
		fmt.Println(err)
	}
	return err
}

func (self *Tube) Dispatch(conn *net.TCPConn) {
	conn.SetKeepAlive(true)
	conn.SetKeepAlivePeriod(time.Second * 60)
	self.fromConn = conn
	self.Recieve()
}

func (self *Tube) ServerDispatch(conn *net.TCPConn) {
	conn.SetKeepAlive(true)
	conn.SetKeepAlivePeriod(time.Second * 60)
	self.fromConn = conn
	self.ServerRecieve()
}
