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
		downloadAddrArg = "0.0.0.0:18080"
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
				upstream, err := net.Dial("tcp", upstreamAddrArg)
				if err != nil {
					fmt.Println("ERROR:", err)
					return
				}
				_ = upstream

				go pipe(conn, upstream, ">>>")
				go pipe(upstream, conn, "<<<")
			}(conn)
		}
	}()

	<-ctx.Done()
}

func pipe(src net.Conn, dst net.Conn, dir string) {
	for {
		buffer := make([]byte, 1024)
		readLen, err := src.Read(buffer)
		if err != nil {
			break
		}

		buffer = buffer[:readLen]

		_, err = dst.Write(buffer)
		if err != nil {
			break
		}

		fmt.Println(dir, string(buffer))
	}
}
