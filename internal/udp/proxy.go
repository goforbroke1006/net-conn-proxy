package udp

import (
	"bufio"
	"fmt"
	"github.com/goforbroke1006/net-conn-proxy/internal/common"
	"net"
)

func InitUDPProxy(downstreamConn *net.UDPConn, upstreamAddr string, bufferSizeArg uint64) {
	for {
		buffer := make([]byte, bufferSizeArg)

		// downstream >>> buffer
		readLen, clientAddr, readErr := downstreamConn.ReadFromUDP(buffer)
		if readErr != nil {
			continue
		}

		go exchange(downstreamConn, clientAddr, buffer[:readLen], upstreamAddr, bufferSizeArg)
	}
}

func exchange(
	downstreamConn *net.UDPConn, clientAddr *net.UDPAddr,
	payload []byte, upstreamAddr string, bufferSizeArg uint64,
) {
	upstreamConn, err := net.Dial("udp", upstreamAddr)
	if err != nil {
		return
	}
	defer func() { _ = upstreamConn.Close() }()

	// buffer >>> upstream
	writeLen, err := upstreamConn.Write(payload)
	if err != nil {
		return
	}

	fmt.Printf("%s (%d) >>> %s (%d)\n%s\n%s\n",
		clientAddr.String(), len(payload),
		upstreamConn.RemoteAddr().String(), writeLen,
		common.GetPrettyHexString(payload),
		string(payload),
	)

	// buffer <<< upstream
	response := make([]byte, bufferSizeArg)
	respLen, respErr := bufio.NewReader(upstreamConn).Read(response)
	if respErr != nil {
		return
	}

	// downstream <<< buffer
	writeLen, writeErr := downstreamConn.WriteToUDP(response[:respLen], clientAddr)
	if writeErr != nil {
		return
	}

	fmt.Printf("%s (%d) <<< %s (%d)\n%s\n%s\n",
		clientAddr.String(), respLen,
		upstreamConn.RemoteAddr().String(), writeLen,
		common.GetPrettyHexString(response[:respLen]),
		string(response[:respLen]),
	)
}
