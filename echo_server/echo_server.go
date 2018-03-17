package main

import (
	"fmt"
	"net"
	"os"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Println(os.Args[0], "<port>")
		return
	}

	fmt.Println("Listening on port: ", os.Args[1])

	ln, err := net.Listen("tcp", ":"+os.Args[1])
	if err != nil {
		fmt.Println("Err listening: ", err.Error())
		return
	}

	for {
		conn, err := ln.Accept()

		if err != nil {
			fmt.Println("Err accepting: ", err.Error())
			return
		}

		go processConn(conn)
	}
}

func processConn(conn net.Conn) {
	for {
		b := make([]byte, 4096)

		n, err := conn.Read(b)
		if err != nil {
			fmt.Println("Err reading: ", err.Error())
			return
		}

		n, err = conn.Write(b[:n])
		if err != nil {
			fmt.Println("Err writing: ", err.Error())
			return
		}
	}
}
