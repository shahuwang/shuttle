package ferry

import (
	"net"
)

type Tunnel struct {
	net.Conn
	priority int
	index    int
}

func (self *Tunnel) Send(data []byte) (n int, err error) {
	// 加密一下发送
	return self.Write(data)
}

func (self *Tunnel) Read(buffer []byte) (n int, err error) {
	// 解密
	n, err = self.Conn.Read(buffer)
	return
}

type TunnelHeap []*Tunnel

func (self TunnelHeap) Len() int {
	return len(self)
}

func (self TunnelHeap) Less(i, j int) bool {
	return self[i].priority < self[j].priority
}

func (self TunnelHeap) Swap(i, j int) {
	self[i], self[j] = self[j], self[i]
	self[i].index = i
	self[j].index = j
}

func (self *TunnelHeap) Push(t interface{}) {
	item := t.(*Tunnel)
	n := len(*self)
	item.index = n
	*self = append(*self, item)
}

func (self *TunnelHeap) Pop() interface{} {
	old := *self
	n := len(old)
	item := old[n-1]
	item.index = -1
	*self = old[0 : n-1]
	return item
}
