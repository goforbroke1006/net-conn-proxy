package tcp

import (
	"fmt"
	"net"
	"sync"

	"go.uber.org/zap"

	"github.com/goforbroke1006/net-conn-proxy/internal/common"
)

func InitTCPProxy(client net.Conn, upstreamAddr string, bufSize uint64) {
	upstream, err := net.Dial("tcp", upstreamAddr)
	if err != nil {
		zap.L().Error("dial downstream", zap.Error(err))
		return
	}

	pipesWg := &sync.WaitGroup{}
	pipesWg.Add(2)
	go pipe(client, upstream, bufSize, pipesWg)
	go pipe(upstream, client, bufSize, pipesWg)
	pipesWg.Wait()

	zap.L().Info("disconnection",
		zap.String("d", client.RemoteAddr().String()),
		zap.String("u", upstream.RemoteAddr().String()))
}

func pipe(src net.Conn, dst net.Conn, bufSize uint64, wg *sync.WaitGroup) {
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
			writeLen, writeErr = dst.Write(buffer[:readLen])

			fmt.Printf("%s (%d) >>> %s (%d)\n%s\n%s\n",
				src.RemoteAddr().String(), readLen,
				dst.RemoteAddr().String(), writeLen,
				common.GetPrettyHexString(buffer[:readLen]),
				string(buffer[:readLen]),
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
