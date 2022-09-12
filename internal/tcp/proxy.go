package tcp

import (
	"fmt"
	"net"
	"strings"
	"sync"

	"github.com/goforbroke1006/net-conn-proxy/internal/common"
)

func InitTCPProxy(client net.Conn, proto, downstreamAddrArg, upstreamAddr string, bufSize uint64) {
	upstream, err := net.Dial(proto, upstreamAddr)
	if err != nil {
		fmt.Println("ERROR:", err)
		return
	}

	var (
		downstreamHost = common.GetHostFromAddr(downstreamAddrArg)
		upstreamHost   = common.GetHostFromAddr(upstreamAddr)
	)

	// dumb mechanism to hide the fact of using a proxy
	// it uses replacing addresses in non-encrypted traffic
	var (
		replacementD2U map[string]string
		replacementU2D map[string]string
	)
	replacementD2U = map[string]string{
		downstreamAddrArg: upstreamAddr,
		downstreamHost:    upstreamHost,
	}
	replacementU2D = make(map[string]string, len(replacementD2U))
	for k, v := range replacementD2U {
		replacementU2D[v] = k
	}

	pipesWg := &sync.WaitGroup{}
	pipesWg.Add(2)
	go pipe(client, upstream, bufSize, pipesWg, replacementD2U)
	go pipe(upstream, client, bufSize, pipesWg, replacementU2D)
	pipesWg.Wait()

	fmt.Printf("disconnection %s <<<>>> %s\n", client.RemoteAddr().String(), upstream.RemoteAddr().String())
}

func pipe(src net.Conn, dst net.Conn, bufSize uint64, wg *sync.WaitGroup, replacement map[string]string) {
	var (
		readLen  int
		readErr  error
		writeLen int
		writeErr error
	)

ExchangeLoop:
	for {
		buffer := make([]byte, bufSize)
		readLen, readErr = src.Read(buffer)

		if readLen > 0 {
			asStr := string(buffer[:readLen])
			for k, v := range replacement {
				asStr = strings.ReplaceAll(asStr, k, v)
			}
			asBytes := []byte(asStr)

			writeLen, writeErr = dst.Write(asBytes)

			fmt.Printf("%s (%d) >>> %s (%d)\n%s\n%s\n",
				src.RemoteAddr().String(), readLen,
				dst.RemoteAddr().String(), writeLen,
				common.GetPrettyHexString(asBytes),
				string(asBytes),
			)
		}

		if readErr != nil {
			break ExchangeLoop
		}

		if writeErr != nil {
			break ExchangeLoop
		}
	}

	_ = src.Close()
	_ = dst.Close()

	wg.Done()
}
