package http

import (
	"fmt"
	"net"
	"strings"
	"sync"

	"go.uber.org/zap"

	"github.com/goforbroke1006/net-conn-proxy/internal/common"
)

func InitHTTPProxy(client net.Conn, downstreamAddrArg, originAddr string, bufSize uint64) {
	upstream, err := net.Dial("tcp", originAddr)
	if err != nil {
		zap.L().Error("dial downstream", zap.Error(err))
		return
	}

	var (
		downstreamHost = common.GetHostFromAddr(downstreamAddrArg)
		upstreamHost   = common.GetHostFromAddr(originAddr)
	)

	// dumb mechanism to hide the fact of using a proxy
	// it uses replacing addresses in non-encrypted traffic
	var (
		replacementD2U map[string]string
		replacementU2D map[string]string
	)
	replacementD2U = map[string]string{
		downstreamAddrArg: originAddr,
		downstreamHost:    upstreamHost,
		"127.0.0.1:8080":  originAddr, // FIXME: dirty workaround for demo example, need better solution for replacing Host in request/response body
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

	zap.L().Info("disconnection",
		zap.String("d", client.RemoteAddr().String()),
		zap.String("u", upstream.RemoteAddr().String()))
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
			zap.L().Error("read from", zap.Error(readErr), zap.String("ip", src.RemoteAddr().String()))
			break ExchangeLoop
		}

		if writeErr != nil {
			zap.L().Error("write to", zap.Error(readErr), zap.String("ip", dst.RemoteAddr().String()))
			break ExchangeLoop
		}
	}

	_ = src.Close()
	_ = dst.Close()

	wg.Done()
}
