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
	clockInterval = 1
)

var (
	mux  = sync.Mutex{}
	pMap = make(PlayerMap)
)

func verifySender(mh MessageHeader, raddr net.Addr) bool {
	storedConnection := pMap.getConnection(mh)

	if storedConnection == nil {
		return false
	} else if storedConnection.(*net.UDPAddr).String() != raddr.(*net.UDPAddr).String() {
		logger.Info(fmt.Sprintf("Player=%v shadow has untracked origin, recorded=%v, received=%v",
			mh.PlayerID, storedConnection, raddr))
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

func handleShadowMessage(msg []byte, pc net.PacketConn, raddr net.Addr) {
	mh := MessageHeader{}
	err := mh.Unmarshall(msg)

	if err == nil {
		if !verifySender(mh, raddr) {
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

func shadowServer() (pc net.PacketConn, err error) {
	pc, err = net.ListenPacket("udp4", fmt.Sprintf(":%d", port))
	logger.Info(fmt.Sprintf("Listening UDP connections on port %d", port))
	if err != nil {
		logger.Error("Error starting UDP shadow server", err)
		return
	}

	buffer := make([]byte, maxBufferSize)

	go func() {
		for {
			n, raddr, err := pc.ReadFrom(buffer)
			if err != nil {
				logger.Error(fmt.Sprintf("Error receiving UDP dgram from %v", raddr), err)
			}
			payloadSize := n - int(unsafe.Sizeof(MessageHeader{}))
			if payloadSize < helloSize {
				logger.Debug("Dgram skipped, payload len=", payloadSize)
				// skip
			} else if payloadSize == helloSize {
				handleHelloMessage(buffer, raddr)
			} else {
				handleShadowMessage(buffer, pc, raddr)
			}
		}
	}()

	return
}

func sessionClock() {
	tick := time.Tick(clockInterval * time.Second)
	for {
		select {
		case <-tick:
			pMap.updateSessions()
		default:
			time.Sleep(clockInterval * time.Second)
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
