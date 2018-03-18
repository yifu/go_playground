package main

import (
	"fmt"
	"net"
)

func main() {
	ln, err := net.Listen("tcp", ":55555")

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

		go processConn(upStreamPeer)
	}
}

func processConn(upStreamPeer net.Conn) {
	downStreamPeer, err := net.Dial("tcp", "localhost:60000")
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
