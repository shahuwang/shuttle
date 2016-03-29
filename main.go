package main

import (
	"fmt"
	"github.com/shahuwang/shuttle/ferry"
	"net"
)

func main() {
	addr, _ := net.ResolveTCPAddr("tcp", "localhost:8088")
	ln, err := net.ListenTCP("tcp", addr)
	if err != nil {
		panic("proxy listen failed")
	}
	go func() {
		for {
			fmt.Println("proxy recieve from server")
			conn, _ := ln.AcceptTCP()
			buffer := ferry.BufferPool()
			for {
				n, err := conn.Read(buffer)
				if err != nil {
					fmt.Println(err)
					fmt.Println("proxy err")
					break
				}
				fmt.Printf("read %d byte proxy", n)
			}
			//conn.Write(buffer)
			//conn.Close()
		}
	}()
	server := ferry.NewServer("localhost:8087", "localhost:8088")
	server.Start()
	fmt.Println("client start")
	client := ferry.NewClient("localhost:8086", "localhost:8087")
	client.Start()
}
