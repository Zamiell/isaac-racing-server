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
	CONN *net.UDPConn
	TTL  uint
}

func (pc *PlayerConn) finalizeConnection(playerId uint32) {
	if pc.CONN != nil {
		err := pc.CONN.Close()
		if err != nil {
			log.Error(fmt.Sprintf("Error closing connection player=%v ", playerId), err)
		} else {
			log.Info(fmt.Sprintf("Connection closed for player=%v", playerId))
		}
	}
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

func (p PlayerMap) getConnection(mh MessageHeader) *net.UDPConn {
	if p[mh.RaceId] != nil && p[mh.RaceId][mh.PlayerId] != nil {
		return p[mh.RaceId][mh.PlayerId].CONN
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
			log.Debug(fmt.Sprintf("[--] player=%v ttl=%v", playerId, pConn.TTL))

			if pConn.TTL <= 0 {
				pConn.finalizeConnection(playerId)
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

func (p *PlayerMap) update(mh MessageHeader, addr string) {
	defer mux.Unlock()

	raddr, _ := net.ResolveUDPAddr("udp", addr)

	mux.Lock()

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
		conn, _ := net.DialUDP("udp4", nil, raddr)
		race[mh.PlayerId] = &PlayerConn{conn, sessionTTL}
		log.Info(fmt.Sprintf("Connection created player=%v, dst=%v", mh.PlayerId, raddr))
	}

	// update TTL
	race[mh.PlayerId].TTL = sessionTTL
	log.Debug(fmt.Sprintf("TTL updated race=%v, player=%v", mh.RaceId, mh.PlayerId))
}
