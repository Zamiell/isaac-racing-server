package main

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
	port          = 9001
	clockInterval = 1 * time.Second
)

var (
	mux  = sync.Mutex{}
	pMap = make(PlayerMap)
)

func verifySender(mh MessageHeader, rAddr net.Addr) bool {
	storedConnection := pMap.getConnection(mh)

	if storedConnection == nil {
		return false
	} else if storedConnection.(*net.UDPAddr).String() != rAddr.(*net.UDPAddr).String() {
		logger.Info(fmt.Sprintf("Player=%v shadow has untracked origin, recorded=%v, received=%v",
			mh.PlayerID, storedConnection, rAddr))
		return false
	}
	return true
}

func handleHelloMessage(msg []byte, addr net.Addr) {
	mh := MessageHeader{}
	err := mh.Unmarshall(msg)
	if err == nil {
		pMap.update(mh, &addr)
	}
}

func handleShadowMessage(msg []byte, pc net.PacketConn, rAddr net.Addr) {
	mh := MessageHeader{}
	err := mh.Unmarshall(msg)

	if err == nil {
		if !verifySender(mh, rAddr) {
			return
		}

		opponent := pMap.getOpponent(mh)
		if opponent != nil {
			_, err := pc.WriteTo(msg, *opponent.ADDR)
			if err != nil {
				logger.Error(fmt.Sprintf(
					"Shadow proxy failed, player=%v, msg: %v\ncause: %v", mh.PlayerID, msg, err))
			}
		}
	}
}

func shadowServer() (net.PacketConn, error) {
	pc, err := net.ListenPacket("udp4", fmt.Sprintf(":%d", port))
	logger.Info(fmt.Sprintf("Listening for UDP connections on port %d", port))
	if err != nil {
		logger.Error("Failed to start the UDP shadow server:", err)
		return nil, err
	}

	buffer := make([]byte, maxBufferSize)

	go func() {
		for {
			n, rAddr, err := pc.ReadFrom(buffer)
			if err != nil {
				logger.Error(fmt.Sprintf("Error receiving UDP datagram from %v", rAddr), err)
				continue
			}
			payloadSize := n - int(unsafe.Sizeof(MessageHeader{}))
			if payloadSize < helloSize {
				logger.Debug("Datagram skipped, payload len=", payloadSize)
				continue
			}
			if payloadSize == helloSize {
				handleHelloMessage(buffer, rAddr)
			} else {
				handleShadowMessage(buffer, pc, rAddr)
			}
		}
	}()

	return pc, nil
}

func sessionClock() {
	tick := time.Tick(clockInterval)
	for {
		select {
		case <-tick:
			pMap.updateSessions()
		default:
			time.Sleep(clockInterval)
		}
	}
}

func shadowInit() {
	go sessionClock()

	var pc net.PacketConn
	if v, err := shadowServer(); err != nil {
		logger.Error("Failed to start the shadow server:", err)
		return
	} else {
		pc = v
	}

	if pc != nil {
		if err := pc.Close(); err != nil {
			logger.Error("Error closing shadow server connection:", err)
			return
		}
	}
}
