package ferry

import (
	"container/heap"
	"errors"
	"fmt"
	"net"
	"sync"
	"time"
)

const CONNUM int = 256
const PACKAGE_SIZE int = 1024 * 8

type Iclient interface {
	Listen()
	Start() error
	HandleConn(conn *net.TCPConn) error
}

type Client struct {
	laddr   string // address for listen
	baddr   string // address to request
	tunnels TunnelHeap
	lock    sync.Mutex
}

func (self *Client) Listen() {
	ln, err := net.Listen("tcp", self.laddr)
	if err != nil {
		fmt.Println("listen failed: %v", err)
		panic("!!")
	}
	listener := ln.(*net.TCPListener)
	for {
		conn, err := listener.AcceptTCP()
		if err != nil {
			if opErr, ok := err.(*net.OpError); ok {
				if !opErr.Temporary() {
					break
				}
			}
			continue
		}
		fmt.Println("accept connection")
		conn.SetKeepAlive(true)
		conn.SetKeepAlivePeriod(time.Second * 60)
		go self.HandleConn(conn)
	}
}

func (self *Client) Start() error {
	heap.Init(&self.tunnels)
	size := cap(self.tunnels)
	for i := 0; i < size; i++ {
		go func(index int) {
			item, err := self.createTunnel()
			if err != nil {
				fmt.Println("tunnel %d reconnect failed", index)
				time.Sleep(time.Second * 3)
				return
			}
			fmt.Println("add tunnel %d", index)
			self.addTunnel(item)
		}(i)
	}
	fmt.Println("start listen")
	self.Listen()
	return nil
}

func (self *Client) HandleConn(conn *net.TCPConn) error {
	defer conn.CloseRead()
	tunnel := self.fetchTunnel()
	if tunnel == nil {
		fmt.Println("no tunnel to use")
		return errors.New("no tunnel to use")
	}
	defer self.dropTunnel(tunnel)
	tube := &Tube{
		tunnel:   tunnel,
		fromConn: conn,
	}
	tube.Dispatch(conn)
	return nil
}

func (self *Client) createTunnel() (tunnel *Tunnel, err error) {
	fmt.Println("start create tunnel")
	conn, err := net.DialTimeout("tcp", self.baddr, time.Duration(5)*time.Second)
	if err != nil {
		fmt.Println("dial server timeout")
		return
	}
	fmt.Println("created tunnel")
	return &Tunnel{Conn: conn}, nil
}

func (self *Client) addTunnel(item *Tunnel) {
	self.lock.Lock()
	defer self.lock.Unlock()
	heap.Push(&self.tunnels, item)
}

func (self *Client) fetchTunnel() *Tunnel {
	defer self.lock.Unlock()
	self.lock.Lock()
	if len(self.tunnels) == 0 {
		return nil
	}
	item := self.tunnels[0]
	item.priority += 1
	heap.Fix(&self.tunnels, 0)
	return item
}

func (self *Client) dropTunnel(item *Tunnel) {
	defer self.lock.Unlock()
	self.lock.Lock()
	item.priority -= 1
	heap.Fix(&self.tunnels, item.index)
}

func BufferPool() []byte {
	return make([]byte, PACKAGE_SIZE)
}

func NewClient(laddr, baddr string) *Client {
	client := &Client{
		laddr:   laddr,
		baddr:   baddr,
		tunnels: make(TunnelHeap, 0, 1),
	}
	return client
}
