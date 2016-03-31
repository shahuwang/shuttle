package ferry

import (
	"bufio"
	"github.com/qiniu/log"
	"net"
	"sync"
	"time"
)

type Tube struct {
	fromConn *net.TCPConn
	tunnel   *Tunnel
	wg       sync.WaitGroup
}

func (self *Tube) Recieve() (err error) {
	defer self.wg.Done()
	bufsize := PACKAGE_SIZE * 2
	rd := bufio.NewReaderSize(self.fromConn, bufsize)
	defer self.fromConn.CloseRead()
	for {
		buffer := BufferPool()
		log.Info("start recieve")
		n, err := rd.Read(buffer)
		if err != nil {
			// TODO
			log.Infof(err.Error())
			break
		}
		log.Info("start send to tunnel")
		_, err = self.tunnel.Send(buffer[:n])
		log.Info("sent to tunnel")
		log.Info(err)
		if err != nil {
			log.Infof(err.Error())
			self.tunnel.Send([]byte("EOF"))
			break
		}
	}
	log.Info("client read over")
	return err
}

func (self *Tube) Back() (err error) {
	defer self.wg.Done()
	// client 端接收请求的回应
	defer self.fromConn.CloseWrite()
	for {
		buffer := BufferPool()
		log.Info("start back")
		n, err := self.tunnel.Read(buffer)
		log.Info("end back")
		if err != nil {
			log.Infof(err.Error())
			break
		}
		_, err = self.fromConn.Write(buffer[:n])
		if err != nil {
			log.Infof(err.Error())
			break
		}
	}
	log.Info("client back over+++++")
	return err
}

func (self *Tube) ServerRecieve() (err error) {
	defer self.wg.Done()
	defer self.fromConn.CloseWrite()
	for {
		buffer := BufferPool()
		n, err := self.tunnel.Read(buffer)
		if err != nil {
			// TODO
			log.Infof(err.Error())
			log.Infof("Recieve error")
			break
		}
		_, err = self.fromConn.Write(buffer[:n])
		if err != nil {
			log.Info("server write to proxy err")
			log.Info(err.Error())
			break
		}
	}
	return err

}

func (self *Tube) ServerBack() (err error) {
	defer self.wg.Done()
	// server 端接收请求的回应
	defer self.fromConn.CloseRead()
	for {
		buffer := BufferPool()
		n, err := self.fromConn.Read(buffer)
		if err != nil {
			log.Infof(err.Error())
			break
		}
		_, err = self.tunnel.Write(buffer[:n])
		if err != nil {
			log.Infof(err.Error())
			break
		}
	}
	return err
}

func (self *Tube) Dispatch(conn *net.TCPConn) {
	log.Info("client Dispatch")
	conn.SetKeepAlive(true)
	conn.SetKeepAlivePeriod(time.Second * 60)
	self.fromConn = conn
	self.wg.Add(1)
	go self.Recieve()
	self.wg.Add(1)
	go self.Back()
	self.wg.Wait()
	log.Info("client dispatch over")
}

func (self *Tube) ServerDispatch(conn *net.TCPConn) {
	conn.SetKeepAlive(true)
	conn.SetKeepAlivePeriod(time.Second * 60)
	self.fromConn = conn
	self.wg.Add(1)
	go self.ServerRecieve()
	self.wg.Add(1)
	go self.ServerBack()
	self.wg.Wait()
}
