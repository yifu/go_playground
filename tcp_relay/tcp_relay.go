package main

import (
	"fmt"
	"net"
	"os"
)

func main() {
	if len(os.Args) != 3 {
		fmt.Println(os.Args[0], "[localAddr]:localPort remoteAddr:remotePort")
		return
	}

	localAddrPort := os.Args[1]
	remoteAddrPort := os.Args[2]

	ln, err := net.Listen("tcp", localAddrPort)
	if err != nil {
		fmt.Println("Err listening:", err.Error())
		return
	}

	for {
		upStreamPeer, err := ln.Accept()
		if err != nil {
			fmt.Println("Err accepting:", err.Error())
			return
		}

		go processConn(upStreamPeer, remoteAddrPort)
	}
}

func processConn(upStreamPeer net.Conn, downStreamAddr string) {
	downStreamPeer, err := net.Dial("tcp", downStreamAddr)
	if err != nil {
		fmt.Println("Remote conn:", err)
		return
	}

	go stream(upStreamPeer, downStreamPeer)
	go stream(downStreamPeer, upStreamPeer)
}

func stream(up, down net.Conn) {
	for {
		b := make([]byte, 4096)

		n, err := up.Read(b)
		if err != nil {
			fmt.Println("Err reading:", err.Error())
			err := down.Close()
			if err != nil {
				fmt.Println("Closing down:", err)
			}
			return
		}

		n, err = down.Write(b[:n])
		if err != nil {
			fmt.Println("Err writing:", err.Error())
			err := up.Close()
			if err != nil {
				fmt.Println("Closing up:", err)
			}
			return
		}
	}

}
