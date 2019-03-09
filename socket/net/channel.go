package net

import (
	"bufio"
	"fmt"
	"net"
	"sync"
)

var bufSize = 1024

type Channel struct {
	conn       net.Conn
	codec      ICodec
	reader     *bufio.Reader
	writer     *bufio.Writer
	readBytes  chan []byte
	writeBytes chan []byte
	revBytes   chan []byte
	sendBytes  chan []byte
	connect    bool
}

func NewChannel(conn net.Conn, codec ICodec) Channel {
	channel := Channel{
		conn:       conn,
		codec:      codec,
		reader:     bufio.NewReader(conn),
		writer:     bufio.NewWriter(conn),
		readBytes:  make(chan []byte),
		writeBytes: make(chan []byte),
		revBytes:   make(chan []byte),
		sendBytes:  make(chan []byte),
		connect:    true}
	channel.Listen()
	return channel
}

var mutex = &sync.Mutex{}

func (channel *Channel) Listen() {
	go func() {
		for {
			buf := <-channel.readBytes
			channel.codec.Decode(buf, channel.revBytes)
		}
	}()
	go func() {
		for {
			buf := <-channel.sendBytes
			channel.codec.Encode(buf, channel.writeBytes)
		}
	}()
	go channel.Read()
	go channel.Write()
}

func (channel *Channel) Read() {
	for channel.connect {
		buf := make([]byte, bufSize)
		n, err := channel.reader.Read(buf)
		if err != nil {
			fmt.Println("err : " + err.Error())
			channel.connect = false
			break
		}
		if n == 0 {
			fmt.Println("empty message")
			channel.connect = false
			break
		} else {
			channel.readBytes <- buf[:n]
		}
	}

}

func (channel *Channel) Write() {
	for data := range channel.writeBytes {
		channel.writer.Write(data)
		channel.writer.Flush()
	}
}