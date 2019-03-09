package net

import (
	"encoding/json"
	"fmt"
	"net"
	"os"
	"strconv"
)

type Configuration struct {
	Port  int
	Codec string
}

type Server struct {
}

func (s Server) Start(codec ICodec) {
	file, err := os.Open("conf.json")
	if err != nil {
		fmt.Println("conf.json not found")
		os.Exit(1)
	}
	decoder := json.NewDecoder(file)
	conf := Configuration{}
	if err != decoder.Decode(&conf) {
		fmt.Println("failed to read conf.json")
		os.Exit(1)
	}
	fmt.Println("loaded conf.json")
	defer file.Close()
	listener, err := net.Listen("tcp", ":"+strconv.Itoa(conf.Port))
	fmt.Println("port binding: " + strconv.Itoa(conf.Port))

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer listener.Close()
	fmt.Println("socket server started")

	process := NewProcess(codec)
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println(err)
		} else {
			process.join <- conn
		}
	}
}
