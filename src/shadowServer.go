// The shadow server simply echos incoming UDP packets back to all of the other players in the race

package server

import (
	"fmt"
	"net"
	"sync"
	"time"
	"unsafe"
)

const (
	maxBufferSize = 1024
	helloSize     = 5
	port          = 9113
	purgeInterval = 1 * time.Second
)

var (
	shadowRaces = ShadowRaces{
		mutex: sync.Mutex{},
		races: make(map[uint32]map[uint32]*PlayerUDPConn),
	}
)

func shadowInit() {
	address := fmt.Sprintf(":%d", port)
	packetConn, err := net.ListenPacket("udp4", address)
	if err != nil {
		logger.Fatal("Failed to start the UDP server:", err)
	}
	logger.Info("Listening for UDP connections on port:", port)

	go UDPServerLoop(packetConn)
	go purgeOldSessionsLoop()
}

func UDPServerLoop(packetConn net.PacketConn) {
	buffer := make([]byte, maxBufferSize)

	for {
		n, addr, err := packetConn.ReadFrom(buffer)
		if err != nil {
			logger.Warning("Failed to read UDP datagram from \""+addr.String()+"\":", err)
			continue
		}

		mh := MessageHeader{}
		if err := mh.Unmarshall(buffer); err != nil {
			logger.Warning("Failed to unmarshall a UDP datagram from \""+addr.String()+"\":", err)
			continue
		}

		payloadSize := n - int(unsafe.Sizeof(MessageHeader{}))
		if payloadSize < helloSize {
			// No message should ever be smaller than a hello message
			continue
		} else if payloadSize == helloSize {
			handleHelloMessage(mh, addr)
		} else {
			handleOtherMessage(mh, addr, packetConn, buffer)
		}
	}
}

func handleHelloMessage(mh MessageHeader, addr net.Addr) {
	// Since we have lazy player initialization,
	// updating the TTL will also instantiate the entry in the map for the respective player
	shadowRaces.updatePlayerTTL(mh, addr)
}

func handleOtherMessage(mh MessageHeader, addr net.Addr, pc net.PacketConn, buffer []byte) {
	if !verifySender(mh, addr) {
		return
	}

	shadowRaces.updatePlayerTTL(mh, addr)

	otherPlayerConnections := shadowRaces.getOtherPlayerConnections(mh)
	for _, conn := range otherPlayerConnections {
		_, err := pc.WriteTo(buffer, conn.addr)
		if err != nil {
			logger.Errorf("Failed to send a UDP message to \"%v\": %w", conn.addr.String(), err)
		}
	}
}

func verifySender(mh MessageHeader, addr net.Addr) bool {
	conn := shadowRaces.getConnection(mh)

	if conn == nil {
		return false
	}

	return conn.addr.String() == addr.String()
}

func purgeOldSessionsLoop() {
	tick := time.Tick(purgeInterval) // nolint: staticcheck

	for {
		select {
		case <-tick:
			shadowRaces.purgeOldSessions()
		default:
			time.Sleep(purgeInterval)
		}
	}
}
