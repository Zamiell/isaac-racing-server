package main

import (
	"fmt"
	"github.com/Zamiell/isaac-racing-server/src/log"
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
	playerConnection := pMap.getConnection(mh)
	if playerConnection == nil {
		log.Debug(fmt.Sprintf("No player=%v session found, ignoring", mh.PlayerId))
		return false
	} else if playerConnection.RemoteAddr().String() != raddr.String() {
		log.Info(fmt.Sprintf("Player=%v shadow comes has untracked origin, recorded=%v, received=%v",
			mh.PlayerId, playerConnection.RemoteAddr().String(), raddr.String()))
		return false
	}
	return true
}

func handleHelloMessage(msg []byte, addr net.Addr) {
	log.Debug("Parsing hello message")
	mh := MessageHeader{}
	err := mh.Unmarshall(msg)
	if err == nil {
		pMap.update(mh, addr.String())
	}
}

func handleShadowMessage(msg []byte, raddr net.Addr) {
	mh := MessageHeader{}
	err := mh.Unmarshall(msg)

	if err == nil {
		log.Debug(fmt.Sprintf("Shadow received, player=%v", mh.PlayerId))
		if !verifySender(mh, raddr) {
			return
		}

		opponent := pMap.getOpponent(mh)
		if opponent != nil {
			log.Debug(fmt.Sprintf("Opponen found, player=%v", mh.PlayerId))
			log.Debug(fmt.Sprintf("Proxying shadow, [src=%v]=>[dst=%v]", raddr, opponent.CONN.RemoteAddr()))
			_, err := opponent.CONN.Write(msg)
			if err != nil {
				log.Debug(fmt.Sprintf(
					"Shadow proxy failed, player=%v, msg: %v\ncause: %v", mh.PlayerId, msg, err))
			}
		} else {
			log.Debug(fmt.Sprintf("Missing opponent session, player=%v", mh.PlayerId))
		}
	}
}

func shadowServer(address string) (err error, pc net.PacketConn) {
	pc, err = net.ListenPacket("udp4", address)
	log.Info("Listening UDP connections on", address)
	if err != nil {
		log.Error("Error starting UDP shadow server", err)
		return
	}

	buffer := make([]byte, maxBufferSize)

	go func() {
		for {
			n, raddr, err := pc.ReadFrom(buffer)
			if err != nil {
				log.Error(fmt.Sprintf("Error receiving UDP dgram from %v", raddr.String()), err)
			}
			if n > 0 {
				log.Debug(fmt.Sprintf("Received dgram on shadow service: bytes=%d src=%s\n", n, raddr.String()))
			}
			payloadSize := n - int(unsafe.Sizeof(MessageHeader{}))
			if payloadSize < helloSize {
				log.Debug("Dgram skipped, payload len=", payloadSize)
				// skip
			} else if payloadSize == helloSize {
				handleHelloMessage(buffer, raddr)
			} else {
				handleShadowMessage(buffer, raddr)
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
	errStart, pc := shadowServer(fmt.Sprintf("127.0.0.1:%d", port))
	if errStart != nil {
		log.Error("Exited by: ", errStart)
		if pc != nil {
			errClose := pc.Close()
			if errClose != nil {
				log.Error("Error closing connection", errClose)
			}
		}
	}
}
