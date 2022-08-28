package main

import (
	"context"
	"flag"
	"fmt"
	"net"
	"os/signal"
)

func main() {
	var (
		upstreamAddrArg string
		downloadAddrArg string = "0.0.0.0:18080"
	)

	flag.StringVar(&upstreamAddrArg, "upstream", upstreamAddrArg, "upstream addr")
	flag.StringVar(&downloadAddrArg, "downstream", downloadAddrArg, "downstream addr")
	flag.Parse()

	ctx, stop := signal.NotifyContext(context.Background())
	defer stop()

	listen, err := net.Listen("tcp", downloadAddrArg)
	if err != nil {
		panic(err)
	}

	go func() {
		for {
			conn, err := listen.Accept()
			if err != nil {
				fmt.Println("ERROR:", err)
				continue
			}

			go func(conn net.Conn) {
				downstream, err := net.Dial("tcp", downloadAddrArg)
				if err != nil {
					fmt.Println("ERROR:", err)
					return
				}
				_ = downstream
			}(conn)
		}
	}()

	<-ctx.Done()
}

func pipe(src net.Conn, dst net.Conn) {
	// TODO:
}
