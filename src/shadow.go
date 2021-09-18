// In seeded races, the silhouettes of other races are drawn onto the screen
// This is accomplished via UDP datagrams that are sent to the client, and then to the server

// The "shadow server" declared in this file is simply a UDP listener
// It expects two different kinds of UDP datagrams:
// 1) beacons, for player initialization and for keeping the connection alive
// 2) shadow data, for transmitting the actual shadow positions to the other players

// The shadow server simply echos incoming non-beacon UDP datagrams back to all of the other players
// in the race without any other processing (besides verifying that the server is coming from the
// right IP address)

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
	beaconSize    = len("HELLO")
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
		if payloadSize < beaconSize {
			// No message should ever be smaller than a beacon message
			continue
		} else if payloadSize == beaconSize {
			handleBeaconMessage(mh, addr)
		} else {
			handleOtherMessage(mh, addr, packetConn, buffer)
		}
	}
}

func handleBeaconMessage(mh MessageHeader, addr net.Addr) {
	logger.Debugf("Got beacon from user %d at address: %s", mh.UserID, addr.String())

	// Since we have lazy player initialization,
	// updating the TTL will also instantiate the entry in the map for the respective player
	shadowRaces.updatePlayerTTL(mh, addr)
}

func handleOtherMessage(mh MessageHeader, addr net.Addr, pc net.PacketConn, buffer []byte) {
	logger.Debugf("Got shadow message from user %d at address: %s", mh.UserID, addr.String())

	if !verifySender(mh, addr) {
		return
	}

	otherPlayerConnections := shadowRaces.getOtherPlayerConnections(mh)
	logger.Debug("Number of other connections:", len(otherPlayerConnections))
	for _, conn := range otherPlayerConnections {
		if n, err := pc.WriteTo(buffer, conn.addr); err != nil {
			logger.Errorf("Failed to send a UDP message to \"%v\": %w", conn.addr.String(), err)
		} else {
			logger.Debugf("Sent shadow message to address: %s (%d bytes)", conn.addr.String(), n)
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
