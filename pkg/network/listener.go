package network

import (
	"bufio"
	"net"
	"strings"
)

func ListenAndServer(addr string, r *Router) error {
	parts := strings.Split(addr, "://")
	var (
		networkProtocol = parts[0]
		serveAddress    = parts[1]
	)

	listen, err := net.Listen(networkProtocol, serveAddress)
	if err != nil {
		return err
	}

	for {
		clientConn, err := listen.Accept()
		if err != nil {
			continue
		}

		go func(c net.Conn) {
			reader := bufio.NewReader(c)
			for {
				line, _, err := reader.ReadLine()
				if err != nil {
					break
				}

				fn := r.getMatch(line)
				fn(line, c)
			}
		}(clientConn)
	}
}
