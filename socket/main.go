package main

import (
	"./net"
)

func main() {
	s := net.Server{}
	s.Start(net.NewPacketCodec())
}
