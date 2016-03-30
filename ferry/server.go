package ferry

import (
	"github.com/qiniu/log"
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
		panic("server listen failed")
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
	tunnel := &Tunnel{Conn: conn}
	fromConn, err := net.DialTCP("tcp", nil, self.baddr)
	if err != nil {
		log.Infof(err.Error())
		log.Infof("server handleconn error")
		return
	}
	log.Infof("handle conn server")
	tube := &Tube{
		tunnel:   tunnel,
		fromConn: fromConn,
	}
	tube.ServerDispatch(fromConn)
	fromConn.Close()
}

func NewServer(lr, br string) *Server {
	laddr, _ := net.ResolveTCPAddr("tcp", lr)
	baddr, _ := net.ResolveTCPAddr("tcp", br)
	return &Server{
		laddr: laddr,
		baddr: baddr,
	}
}
