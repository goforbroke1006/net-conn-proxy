package main

import (
	"context"
	"flag"
	"fmt"
	"net"
	"os"
	"os/signal"
	"strings"
	"sync"

	"github.com/goforbroke1006/net-conn-proxy/internal"
)

func main() {
	var (
		protocolArg       string
		downstreamAddrArg = "0.0.0.0:0"
		upstreamAddrArg   string
		bufferSizeArg     uint64 = 2048
	)

	flag.StringVar(&protocolArg, "p", protocolArg, "protocol - tcp or udp")
	flag.StringVar(&downstreamAddrArg, "d", downstreamAddrArg, "downstream addr like 120.0.0.1:8080")
	flag.StringVar(&upstreamAddrArg, "u", upstreamAddrArg, "upstream addr like 8.8.8.8:80")
	flag.Uint64Var(&bufferSizeArg, "bs", bufferSizeArg, "pipe buffer size")
	flag.Parse()

	if len(protocolArg) == 0 {
		fmt.Println("ERROR: -p argument required")
		os.Exit(1)
	}
	if protocolArg != "tcp" && protocolArg != "udp" {
		fmt.Println("ERROR: -p argument allow 'tcp' or 'upd' value")
		os.Exit(1)
	}
	if len(downstreamAddrArg) == 0 {
		fmt.Println("ERROR: -d argument required")
		os.Exit(1)
	}
	if len(upstreamAddrArg) == 0 {
		fmt.Println("ERROR: -u argument required")
		os.Exit(1)
	}

	fmt.Println("proxying", upstreamAddrArg, "on", protocolArg, "buffer size", bufferSizeArg)

	ctx, stop := signal.NotifyContext(context.Background())
	defer stop()

	listener, err := net.Listen(protocolArg, downstreamAddrArg)
	if err != nil {
		panic(err)
	}

	go func() {
		for {
			clientConn, err := listener.Accept()
			if err != nil {
				fmt.Println("ERROR:", err)
				continue
			}
			go proxy(clientConn, protocolArg, downstreamAddrArg, upstreamAddrArg, bufferSizeArg)
		}
	}()

	<-ctx.Done()
}

func proxy(client net.Conn, proto, downstreamAddrArg, upstreamAddr string, bufSize uint64) {
	upstream, err := net.Dial(proto, upstreamAddr)
	if err != nil {
		fmt.Println("ERROR:", err)
		return
	}

	downstreamHost := internal.GetHostFromAddr(downstreamAddrArg)
	upstreamHost := internal.GetHostFromAddr(upstreamAddr)

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
				internal.GetPrettyHexString(asBytes),
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
