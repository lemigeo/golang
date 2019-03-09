package net

import (
	"bytes"
	"encoding/binary"
)

const (
	None = iota
	Connect
	Hello
)

type IPacket interface {
	ID() int
	ToBytes() []byte
	ToPacket() IPacket
}

type ConnectPacket struct {
	Table []byte
}

func (p ConnectPacket) ID() int {
	return Connect
}

func (p ConnectPacket) ToBytes() []byte {
	buffer := new(bytes.Buffer)
	binary.Write(buffer, binary.LittleEndian, int32(p.ID()))
	binary.Write(buffer, binary.LittleEndian, p.Table)
	return buffer.Bytes()
}

type HelloPacket struct {
	Msg string
}

func (p HelloPacket) ID() int {
	return Hello
}

func (p HelloPacket) ToBytes() []byte {
	buffer := new(bytes.Buffer)
	binary.Write(buffer, binary.LittleEndian, int32(p.ID()))
	binary.Write(buffer, binary.LittleEndian, int16(len(p.Msg)))
	binary.Write(buffer, binary.LittleEndian, []byte(p.Msg))
	return buffer.Bytes()
}
