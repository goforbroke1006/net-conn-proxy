package main

import (
	"net"

	"github.com/goforbroke1006/proxy/pkg/network"
)

func main() {
	router := &network.Router{}
	router.HandleFunc(`^INFO$`, func(payload []byte, w net.Conn) {
		_, _ = w.Write([]byte("TCP demo server\n"))
		_, _ = w.Write([]byte("Version: v0.0.0\n"))
	})
	router.HandleFunc(`^QUIT$`, func(payload []byte, w net.Conn) {
		_ = w.Close()
	})
	router.HandleFunc(``, func(payload []byte, w net.Conn) {
		_, _ = w.Write([]byte("WRONG COMMAND\n"))
	})

	if err := network.ListenAndServer("tcp://0.0.0.0:7777", router); err != nil {
		panic(err)
	}
}
