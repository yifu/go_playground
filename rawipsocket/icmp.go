package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"log"
	"net"
)

type icmpEcho struct {
	typ, code          uint8
	checksum           uint16
	identifier, seqnum uint16
	data               []byte
}

func (pkt icmpEcho) String() string {
	return fmt.Sprintf("{type: %d, code: %d, checksum: %d, id: %d, seqnum: %d, data: %v}",
		pkt.typ, pkt.code, pkt.checksum, pkt.identifier, pkt.seqnum, pkt.data)
}

func main() {
	netaddr, _ := net.ResolveIPAddr("ip4", "en0")
	conn, err := net.ListenIP("ip4:icmp", netaddr)
	if err != nil {
		log.Fatal(err)
	}

	buf := make([]byte, 1024)
	numRead, _, _ := conn.ReadFrom(buf)
	if numRead <= 2 {
		log.Fatal("Received a pkt too small: ", numRead, " long")
	}

	typ, n := binary.Uvarint(buf[0:2])
	if n == 0 {
		log.Fatal("Reading type field led to error: ", n)
	}

	switch {
	case typ == 0 || typ == 8:
		processIcmpEcho(buf)
	}
}

func processIcmpEcho(buf []byte) {
	var pkt icmpEcho
	r := bytes.NewReader(buf[:])
	binary.Read(r, binary.BigEndian, &pkt.typ)
	binary.Read(r, binary.BigEndian, &pkt.code)
	binary.Read(r, binary.BigEndian, &pkt.checksum)
	binary.Read(r, binary.BigEndian, &pkt.identifier)
	binary.Read(r, binary.BigEndian, &pkt.seqnum)
	binary.Read(r, binary.BigEndian, &pkt.data)

	fmt.Println(pkt)
}
