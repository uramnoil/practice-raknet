package main

import (
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"net"
	"time"
)

func main() {

	var guid = 6537344134170577065

	serverAddr, _ := net.ResolveUDPAddr("udp", ":19132")
	serverConn, _ := net.ListenUDP("udp", serverAddr)
	defer serverConn.Close()

	fmt.Println("Server started on port 19132")

	i := 0

	for {
		packet := make([]byte, 4096)
		_, addr, _ := serverConn.ReadFromUDP(packet)

		if packet[0] == 0x01 {
			clientTime := uint64(binary.BigEndian.Uint64(packet[1:9]))
			clientMagic := packet[9:25]
			clientGuid := binary.BigEndian.Uint64(packet[25:])

			fmt.Printf("Received UNCONNECTED_PING packet(time: %d, magic %s, guid: %d\n", clientTime, hex.EncodeToString(clientMagic), clientGuid)

			serverTime := time.Now()

			magic := []byte {0x00, 0xff, 0xff, 0x00, 0xfe, 0xfe, 0xfe, 0xfe, 0xfd, 0xfd, 0xfd, 0xfd, 0x12, 0x34, 0x56, 0x78}

			serverIdString := fmt.Sprintf("MCPE;UNKO;762;1.19.4;0;100;13253860892328930865;Counter: %d;Survival;1;19132;19133;", i)
			response := make([]byte, 1 + 8 + 8 + 16 + 2 + len(serverIdString) + 1)
			response[0] = 0x1c
			binary.BigEndian.PutUint64(response[1:], uint64(serverTime.UnixMilli()))
			binary.BigEndian.PutUint64(response[9:], uint64(guid))
			copy(response[17:], magic)
			binary.BigEndian.PutUint16(response[33:], uint16(len(serverIdString)))
			copy(response[35:], []byte(serverIdString))
			serverConn.WriteToUDP(response, addr)

			i += 1
		}
	}
}
