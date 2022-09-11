package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/goforbroke1006/net-conn-proxy/internal"
	"net"
	"os"
	"os/signal"
	"strings"
	"sync"
)

func main() {
	var (
		protocolArg       string
		downstreamAddrArg = "0.0.0.0:0"
		upstreamAddrArg   string
	)

	flag.StringVar(&protocolArg, "p", protocolArg, "protocol - tcp or udp")
	flag.StringVar(&downstreamAddrArg, "d", downstreamAddrArg, "downstream addr like 120.0.0.1:8080")
	flag.StringVar(&upstreamAddrArg, "u", upstreamAddrArg, "upstream addr like 8.8.8.8:80")
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
			go proxy(clientConn, protocolArg, downstreamAddrArg, upstreamAddrArg)
		}
	}()

	<-ctx.Done()
}

func proxy(client net.Conn, proto, downstreamAddrArg, upstreamAddr string) {
	upstream, err := net.Dial(proto, upstreamAddr)
	if err != nil {
		fmt.Println("ERROR:", err)
		return
	}

	downstreamHost := internal.GetHostFromAddr(downstreamAddrArg)
	upstreamHost := internal.GetHostFromAddr(upstreamAddr)

	pipesWg := &sync.WaitGroup{}
	pipesWg.Add(2)
	go pipe(client, upstream, pipesWg, map[string]string{downstreamAddrArg: upstreamAddr, downstreamHost: upstreamHost})
	go pipe(upstream, client, pipesWg, map[string]string{upstreamAddr: downstreamAddrArg, upstreamHost: downstreamHost})
	pipesWg.Wait()

	fmt.Printf("disconnection %s <<<>>> %s\n", client.RemoteAddr().String(), upstream.RemoteAddr().String())
}

func pipe(src net.Conn, dst net.Conn, wg *sync.WaitGroup, replacement map[string]string) {
	var (
		readLen  int
		readErr  error
		writeLen int
		writeErr error
	)

ExchangeLoop:
	for {
		buffer := make([]byte, 1024)
		readLen, readErr = src.Read(buffer)

		if readLen > 0 {
			asStr := string(buffer[:readLen])
			for k, v := range replacement {
				asStr = strings.ReplaceAll(asStr, k, v)
			}
			asBytes := []byte(asStr)

			writeLen, writeErr = dst.Write(asBytes)

			fmt.Printf("%s (%d) >>> %s (%d)\n%s\n",
				src.RemoteAddr().String(), readLen,
				dst.RemoteAddr().String(), writeLen,
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
