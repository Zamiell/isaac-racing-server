package server

import (
	"bytes"
	"encoding/binary"
)

type MessageHeader struct {
	RaceID uint32
	UserID uint32
}

func (mh *MessageHeader) Unmarshall(b []byte) error {
	reader := bytes.NewReader(b)
	return binary.Read(reader, binary.LittleEndian, mh)
}
