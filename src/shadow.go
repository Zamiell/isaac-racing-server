package server

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"net"
)

const (
	sessionTTL = 300
)

type PlayerMap map[uint32]map[uint32]*PlayerConn

type PlayerConn struct {
	ADDR *net.Addr
	TTL  uint
}

type MessageHeader struct {
	RaceID   uint32
	PlayerID uint32
}

func (m *MessageHeader) Unmarshall(b []byte) (err error) {
	reader := bytes.NewReader(b)
	err = binary.Read(reader, binary.LittleEndian, m)
	if err != nil {
		logger.Error(fmt.Sprintf("Failed to unmarshall message: %v", b), err)
	}
	return
}

func (p PlayerMap) getConnection(mh MessageHeader) net.Addr {
	if p[mh.RaceID] != nil && p[mh.RaceID][mh.PlayerID] != nil {
		return *p[mh.RaceID][mh.PlayerID].ADDR
	}
	return nil
}

func (p PlayerMap) getOpponent(mh MessageHeader) (pConn *PlayerConn) {
	race := p[mh.RaceID]
	if race == nil {
		return
	}
	for pID, conn := range race {
		if pID != mh.PlayerID {
			pConn = conn
		}
	}
	return
}

func (p *PlayerMap) updateSessions() {
	mux.Lock()
	defer mux.Unlock()

	for raceID, race := range *p {
		for playerID, pConn := range race {
			if pConn == nil {
				continue
			}

			pConn.TTL--
			if pConn.TTL <= 0 {
				delete(race, playerID)
				logger.Debug(fmt.Sprintf("Removing player=%v from race=%v due to timeout", playerID, raceID))

				if len(race) < 1 {
					delete(*p, raceID)
					logger.Debug(fmt.Sprintf("Record removed race=%v", raceID))
				}
			}
		}
	}
}

func (p *PlayerMap) update(mh MessageHeader, addr *net.Addr) {
	mux.Lock()
	defer mux.Unlock()

	// lazy-init race
	if (*p)[mh.RaceID] == nil {
		(*p)[mh.RaceID] = make(map[uint32]*PlayerConn)
		logger.Debug(fmt.Sprintf("Record created race=%v ", mh.RaceID))
	}

	race := (*p)[mh.RaceID]

	// lazy-init player connection
	if race[mh.PlayerID] == nil {
		if len(race) > 1 {
			logger.Info(fmt.Sprintf("Player=%v attempted to join race=%v with two players", mh.PlayerID, mh.RaceID))
			return
		}
		race[mh.PlayerID] = &PlayerConn{addr, sessionTTL}
		logger.Info(fmt.Sprintf("Connection created player=%v, dst=%v", mh.PlayerID, *addr))
	}

	// update TTL
	race[mh.PlayerID].TTL = sessionTTL
}
