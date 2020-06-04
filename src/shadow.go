package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"github.com/Zamiell/isaac-racing-server/src/log"
	"net"
)

const (
	sessionTTL = 30
)

type PlayerMap map[uint32]map[uint32]*PlayerConn

type PlayerConn struct {
	ADDR *net.Addr
	TTL  uint
}

type MessageHeader struct {
	RaceId   uint32
	PlayerId uint32
}

func (m *MessageHeader) Unmarshall(b []byte) (err error) {
	reader := bytes.NewReader(b)
	err = binary.Read(reader, binary.LittleEndian, m)
	if err != nil {
		log.Error(fmt.Sprintf("Failed to unmarshall message: %v", b), err)
	}
	return
}

func (p PlayerMap) getConnection(mh MessageHeader) net.Addr {
	if p[mh.RaceId] != nil && p[mh.RaceId][mh.PlayerId] != nil {
		return *p[mh.RaceId][mh.PlayerId].ADDR
	}
	return nil
}

func (p PlayerMap) getOpponent(mh MessageHeader) (pConn *PlayerConn) {
	race := p[mh.RaceId]
	if race == nil {
		return
	}
	for pID, conn := range race {
		if pID != mh.PlayerId {
			pConn = conn
		}
	}
	return
}

func (p *PlayerMap) updateSessions() {
	mux.Lock()
	defer mux.Unlock()

	for raceId, race := range *p {
		for playerId, pConn := range race {
			if pConn == nil {
				continue
			}

			pConn.TTL--
			if pConn.TTL <= 0 {
				delete(race, playerId)
				log.Debug(fmt.Sprintf("Removing player=%v from race=%v due to timeout", playerId, raceId))

				if len(race) < 1 {
					delete(*p, raceId)
					log.Debug(fmt.Sprintf("Record removed race=%v", raceId))
				}
			}
		}
	}
}

func (p *PlayerMap) update(mh MessageHeader, addr *net.Addr) {
	mux.Lock()
	defer mux.Unlock()

	// lazy-init race
	if (*p)[mh.RaceId] == nil {
		(*p)[mh.RaceId] = make(map[uint32]*PlayerConn)
		log.Debug(fmt.Sprintf("Record created race=%v ", mh.RaceId))
	}

	race := (*p)[mh.RaceId]

	// lazy-init player connection
	if race[mh.PlayerId] == nil {
		if len(race) > 1 {
			log.Info(fmt.Sprintf("Player=%v attempted to join race=%v with two players", mh.PlayerId, mh.RaceId))
			return
		}
		race[mh.PlayerId] = &PlayerConn{addr, sessionTTL}
		log.Info(fmt.Sprintf("Connection created player=%v, dst=%v", mh.PlayerId, *addr))
	}

	// update TTL
	race[mh.PlayerId].TTL = sessionTTL
}
