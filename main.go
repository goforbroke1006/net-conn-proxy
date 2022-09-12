package main

import (
	"context"
	"flag"
	"fmt"
	"net"
	"os"
	"os/signal"

	"github.com/goforbroke1006/net-conn-proxy/internal/common"
	"github.com/goforbroke1006/net-conn-proxy/internal/tcp"
	"github.com/goforbroke1006/net-conn-proxy/internal/udp"
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

	if protocolArg == "tcp" {
		listener, err := net.Listen("tcp", downstreamAddrArg)
		if err != nil {
			panic(err)
		}
		defer func() { _ = listener.Close() }()

		go func() {
			for {
				clientConn, err := listener.Accept()
				if err != nil {
					fmt.Println("ERROR:", err)
					continue
				}
				go tcp.InitTCPProxy(clientConn, protocolArg, downstreamAddrArg, upstreamAddrArg, bufferSizeArg)
			}
		}()
	}

	if protocolArg == "udp" {
		addr := &net.UDPAddr{
			IP:   net.ParseIP(common.GetHostFromAddr(downstreamAddrArg)),
			Port: int(common.GetPortFromAddr(downstreamAddrArg)),
		}
		downstreamConn, err := net.ListenUDP("udp", addr)
		if err != nil {
			panic(err)
		}
		defer func() { _ = downstreamConn.Close() }()

		go udp.InitUDPProxy(downstreamConn, upstreamAddrArg, bufferSizeArg)
	}

	<-ctx.Done()
	fmt.Println("bye...")
}
