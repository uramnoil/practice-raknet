package main

import (
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"net"
	"time"
)

func main() {

	guid := 6537344134170577065
	magic := []byte {0x00, 0xff, 0xff, 0x00, 0xfe, 0xfe, 0xfe, 0xfe, 0xfd, 0xfd, 0xfd, 0xfd, 0x12, 0x34, 0x56, 0x78}

	serverAddr, _ := net.ResolveUDPAddr("udp", ":19132")
	serverConn, _ := net.ListenUDP("udp", serverAddr)
	defer serverConn.Close()

	fmt.Println("Server started on port 19132")

	i := 0

	for {
		packet := make([]byte, 4096)
		_, addr, _ := serverConn.ReadFromUDP(packet)

		if packet[0] == 0x01 {

			// <-- Request -->

			clientTime := uint64(binary.BigEndian.Uint64(packet[1:9]))
			clientMagic := packet[9:25]
			clientGuid := binary.BigEndian.Uint64(packet[25:])

			fmt.Printf("Received UNCONNECTED_PING packet(time: %d, magic %s, guid: %d\n", clientTime, hex.EncodeToString(clientMagic), clientGuid)

			// <-- Response -->

			serverTime := time.Now()
			serverIdString := fmt.Sprintf("MCPE;UNKO;762;1.19.4;0;100;13253860892328930865;Counter: %d;Survival;1;19132;19133;", i)

			response := []byte{}

			// Packet ID
			response = append(response, 0x1c)

			// Timestamp
			response = binary.BigEndian.AppendUint64(response, uint64(serverTime.UnixMilli()))

			// GUID
			response = binary.BigEndian.AppendUint64(response, uint64(guid))

			// MAGIC
			response = append(response, magic...)

			// Length of Server ID String
			response = binary.BigEndian.AppendUint16(response, uint16(len(serverIdString)))

			/// Server ID String
			response = append(response, []byte(serverIdString)...)

			serverConn.WriteToUDP(response, addr)

			i += 1
		}
	}
}
