package net

import (
	"bytes"
	"crypto/rand"
	"encoding/binary"
	"fmt"
	"net"
)

type Process struct {
	channels   []Channel
	sessions   map[Channel]int64
	join       chan net.Conn
	do         chan []byte
	send       chan []byte
	readBytes  chan []byte
	writeBytes chan []byte
	codec      ICodec
}

func NewProcess(codec ICodec) Process {
	process := Process{
		channels:   make([]Channel, 0),
		sessions:   make(map[Channel]int64, 0),
		join:       make(chan net.Conn),
		do:         make(chan []byte),
		send:       make(chan []byte),
		readBytes:  make(chan []byte),
		writeBytes: make(chan []byte),
		codec:      codec}
	process.Run()
	return process
}

func (process *Process) Run() {
	go func() {
		for {
			select {
			case conn := <-process.join:
				process.Connected(conn)
			case data := <-process.do:
				process.Received(data)
			}
		}
	}()
}

func (process *Process) Connected(conn net.Conn) {
	fmt.Println("connected a new user")
	codec := process.codec.Clone()
	channel := NewChannel(conn, codec)
	process.channels = append(process.channels, channel)
	go func() {
		for {
			process.do <- <-channel.revBytes
		}
	}()
	go func() {
		for {
			channel.sendBytes <- <-process.send
		}
	}()
	table := make([]byte, 30)
	rand.Read(table)
	packet := &ConnectPacket{Table: table}
	process.send <- packet.ToBytes()
	channel.codec.SetBytes(table)
}

func (process *Process) Received(data []byte) {
	fmt.Println("recevied packet")
	var id int32
	buffer := bytes.NewBuffer(data)
	err := binary.Read(buffer, binary.LittleEndian, &id)
	if err != nil {
		fmt.Println(err)
	}
	//byte parsing and do business
	//result set below
	switch int(id) {
	case Hello:
		var size int16
		binary.Read(buffer, binary.LittleEndian, &size)
		msg := make([]byte, size)
		binary.Read(buffer, binary.LittleEndian, msg)
		fmt.Println(string(msg))
		process.send <- data
		break
	}
}
