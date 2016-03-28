package ferry

import (
	"net"
)

type Server struct {
	laddr *net.TCPAddr // client数据
	baddr *net.TCPAddr
}

func (self *Server) Start() {
	go self.Listen()
}

func (self *Server) Listen() {
	ln, err := net.ListenTCP("tcp", self.laddr)
	if err != nil {
		panic("listen failed:%v", err)
	}
	for {
		conn, err := ln.AcceptTCP()
		if err != nil {
			if opErr, ok := err.(*net.OpError); ok {
				if !opErr.Temporary() {
					break
				}
			}
			continue
		}
		go self.HandleConn(conn)
	}
}

func (self *Server) HandleConn(conn *net.TCPConn) {
	defer conn.Close()

}
