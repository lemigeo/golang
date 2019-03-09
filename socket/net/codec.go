package net

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"math"
)

type ICodec interface {
	SetBytes(bytes []byte)
	Encode(bytes []byte, writeBytes chan []byte)
	Decode(readBytes []byte, bytes chan []byte)
	Clone() ICodec
}

type PacketCodec struct {
	isReadHeader bool
	headerSize   int32
	headerBuf    []byte
	packetBuf    []byte
	totalLength  int32
	leftByteSize int32
	headerOffset int32
	packetOffset int32
	xorBytes     []byte
}

func NewPacketCodec() PacketCodec {
	return PacketCodec{
		isReadHeader: false,
		headerSize:   4,
		headerBuf:    make([]byte, 4),
		packetBuf:    nil,
		totalLength:  0,
		leftByteSize: 0,
		headerOffset: 0,
		xorBytes:     nil}
}

func (codec PacketCodec) SetBytes(bytes []byte) {
	codec.xorBytes = bytes
}

func (codec PacketCodec) Encode(buffer []byte, writeBytes chan []byte) {
	if codec.xorBytes != nil {
		buffer = codec.xor(buffer)
	}
	contentLength := int32(len(buffer))
	packet := make([]byte, contentLength+codec.headerSize)
	buf := new(bytes.Buffer)
	err := binary.Write(buf, binary.LittleEndian, &contentLength)
	if err != nil {
		fmt.Println("byte error")
	} else {
		copy(packet[0:], buf.Bytes())
		copy(packet[4:], buffer)
		writeBytes <- packet
	}
}

func (codec PacketCodec) xor(data []byte) []byte {
	buffer := make([]byte, len(data))
	tableIdx := 0
	for i := 0; i < len(buffer); i++ {
		if tableIdx >= len(codec.xorBytes) {
			tableIdx = 0
		}
		buffer[i] = data[i] ^ codec.xorBytes[tableIdx]
	}
	return buffer
}

func (codec PacketCodec) Decode(readBytes []byte, revBytes chan []byte) {
	offset := int32(0)
	size := int32(len(readBytes))
	for size > 0 {
		if !codec.isReadHeader {
			headerByteSize := int32(math.Min(float64(codec.headerSize), float64(size)))
			copy(codec.headerBuf[codec.headerOffset:], readBytes[offset:])
			codec.headerOffset += headerByteSize
			if headerByteSize < codec.headerSize {
				break
			}
			if codec.xorBytes != nil {
				codec.headerBuf = codec.xor(codec.headerBuf)
			}
			buf := bytes.NewBuffer(codec.headerBuf)
			binary.Read(buf, binary.LittleEndian, &codec.totalLength)

			size -= headerByteSize
			offset += headerByteSize
			codec.leftByteSize = codec.totalLength
			codec.packetBuf = make([]byte, codec.totalLength)
			codec.isReadHeader = true
		}
		readByteSize := int32(math.Min(float64(codec.leftByteSize), float64(size)))
		copy(codec.packetBuf[codec.totalLength-codec.leftByteSize:], readBytes[offset:])
		codec.leftByteSize -= readByteSize
		if codec.leftByteSize == 0 {
			if codec.xorBytes != nil {
				codec.packetBuf = codec.xor(codec.packetBuf)
			}
			revBytes <- codec.packetBuf[:]
			codec.clear()
		}
		offset += readByteSize
		size -= readByteSize
	}
}

func (codec PacketCodec) decrypt(data []byte) []byte {
	return nil
}

func (codec PacketCodec) Clone() ICodec {
	clone := NewPacketCodec()
	if codec.xorBytes != nil {
		clone.SetBytes(codec.xorBytes)
	}
	return clone
}

func (codec PacketCodec) clear() {
	codec.leftByteSize = 0
	codec.headerOffset = 0
	codec.isReadHeader = false
	codec.totalLength = 0
	codec.packetBuf = nil
}

type JsonCodec struct {
	isReadHeader bool
	headerSize   int32
	headerBuf    []byte
	packetBuf    []byte
	totalLength  int32
	leftByteSize int32
	headerOffset int32
	packetOffset int32
}

func NewJsonCodec() JsonCodec {
	return JsonCodec{
		isReadHeader: false,
		headerSize:   4,
		headerBuf:    make([]byte, 4),
		packetBuf:    nil,
		totalLength:  0,
		leftByteSize: 0,
		headerOffset: 0}
}

func (codec JsonCodec) SetBytes(bytes []byte) {
}

func (codec JsonCodec) Encode(buffer []byte, writeBytes chan []byte) {
	contentLength := int32(len(buffer))
	packet := make([]byte, contentLength+codec.headerSize)
	buf := new(bytes.Buffer)
	err := binary.Write(buf, binary.BigEndian, &contentLength)
	if err != nil {
		fmt.Println("byte error")
	} else {
		copy(packet[0:], buf.Bytes())
		copy(packet[4:], buffer)
		writeBytes <- packet
	}
}

func (codec JsonCodec) Decode(readBytes []byte, revBytes chan []byte) {
	offset := int32(0)
	size := int32(len(readBytes))
	for size > 0 {
		if !codec.isReadHeader {
			headerByteSize := int32(math.Min(float64(codec.headerSize), float64(size)))
			copy(codec.headerBuf[codec.headerOffset:], readBytes[offset:])
			codec.headerOffset += headerByteSize
			if headerByteSize < codec.headerSize {
				break
			}
			buf := bytes.NewBuffer(codec.headerBuf)
			binary.Read(buf, binary.BigEndian, &codec.totalLength)

			size -= headerByteSize
			offset += headerByteSize
			codec.leftByteSize = codec.totalLength
			codec.packetBuf = make([]byte, codec.totalLength)
			codec.isReadHeader = true
		}
		readByteSize := int32(math.Min(float64(codec.leftByteSize), float64(size)))
		copy(codec.packetBuf[codec.totalLength-codec.leftByteSize:], readBytes[offset:])
		codec.leftByteSize -= readByteSize
		if codec.leftByteSize == 0 {
			revBytes <- codec.packetBuf[:]
			codec.clear()
		}
		offset += readByteSize
		size -= readByteSize
	}
}

func (codec JsonCodec) Clone() ICodec {
	return NewJsonCodec()
}

func (codec JsonCodec) clear() {
	codec.leftByteSize = 0
	codec.headerOffset = 0
	codec.isReadHeader = false
	codec.totalLength = 0
	codec.packetBuf = nil
}