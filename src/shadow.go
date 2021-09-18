package server

import (
	"bytes"
	"encoding/binary"
	"net"
	"sync"
)

const (
	UDPSessionTTLSeconds = 5 * 60
)

type MessageHeader struct {
	RaceID   uint32
	PlayerID uint32
}

func (mh *MessageHeader) Unmarshall(b []byte) error {
	reader := bytes.NewReader(b)
	return binary.Read(reader, binary.LittleEndian, mh)
}

type ShadowRaces struct {
	mutex sync.Mutex
	races map[uint32]map[uint32]*PlayerUDPConn
}

type PlayerUDPConn struct {
	addr net.Addr
	TTL  uint
}

func (sr *ShadowRaces) updatePlayerTTL(mh MessageHeader, addr net.Addr) {
	sr.mutex.Lock()
	defer sr.mutex.Unlock()

	// Lazy-init the player map for every race
	players, ok := sr.races[mh.RaceID]
	if !ok {
		players = make(map[uint32]*PlayerUDPConn)
		sr.races[mh.RaceID] = players
	}

	// Lazy-init the player connection
	conn, ok := players[mh.PlayerID]
	if !ok {
		conn = &PlayerUDPConn{addr, UDPSessionTTLSeconds}
		players[mh.PlayerID] = conn
	}

	conn.TTL = UDPSessionTTLSeconds
}

func (sr *ShadowRaces) getConnection(mh MessageHeader) *PlayerUDPConn {
	// We do not need to acquire the mutex if we are just reading values

	players, ok := sr.races[mh.RaceID]
	if !ok {
		return nil
	}

	conn, ok := players[mh.PlayerID]
	if !ok {
		return nil
	}

	return conn
}

func (sr *ShadowRaces) getOtherPlayerConnections(mh MessageHeader) []*PlayerUDPConn {
	// We do not need to acquire the mutex if we are just reading values

	players, ok := sr.races[mh.RaceID]
	if !ok {
		return nil
	}

	otherPlayerConnections := make([]*PlayerUDPConn, 0)
	for playerID, conn := range players {
		if playerID != mh.PlayerID {
			otherPlayerConnections = append(otherPlayerConnections, conn)
		}
	}

	return otherPlayerConnections
}

func (sr *ShadowRaces) purgeOldSessions() {
	sr.mutex.Lock()
	defer sr.mutex.Unlock()

	for raceID, players := range sr.races {
		for playerID, conn := range players {
			if conn == nil {
				continue
			}

			conn.TTL--

			if conn.TTL > 0 {
				continue
			}

			delete(players, playerID)
			if len(players) == 0 {
				delete(sr.races, raceID)
			}
		}
	}
}
