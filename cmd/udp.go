package cmd

import (
	"context"
	"fmt"
	"net"
	"os/signal"

	"github.com/spf13/cobra"
	"go.uber.org/zap"

	"github.com/goforbroke1006/net-conn-proxy/internal/common"
	"github.com/goforbroke1006/net-conn-proxy/internal/component/udp"
)

func NewUDPCmd() *cobra.Command {
	var (
		downstreamAddrArg = "0.0.0.0:0"
		upstreamAddrArg   string
		bufferSizeArg     uint64 = 2048
	)

	cmd := &cobra.Command{
		Use: "udp",
		Run: func(cmd *cobra.Command, args []string) {
			ctx, stop := signal.NotifyContext(context.Background())

			addr := &net.UDPAddr{
				IP:   net.ParseIP(common.GetHostFromAddr(downstreamAddrArg)),
				Port: int(common.GetPortFromAddr(downstreamAddrArg)),
			}
			downstreamConn, err := net.ListenUDP("udp", addr)
			if err != nil {
				zap.L().Fatal("downstream UDP", zap.Error(err))
			}
			defer func() { _ = downstreamConn.Close() }()

			zap.L().Info("TCP proxying",
				zap.String("up", upstreamAddrArg),
				zap.Uint64("buffer", bufferSizeArg))

			proxy := udp.NewProxyDataGram(downstreamConn, upstreamAddrArg, bufferSizeArg)
			go proxy.Run()

			<-ctx.Done()
			defer stop()
			fmt.Println("bye...")
		},
	}

	cmd.PersistentFlags().StringVarP(&downstreamAddrArg, "downstream", "d", downstreamAddrArg,
		"downstream addr like 120.0.0.1:8080")
	cmd.PersistentFlags().StringVarP(&upstreamAddrArg, "upstream", "u", upstreamAddrArg,
		"upstream addr like 8.8.8.8:80")
	cmd.PersistentFlags().Uint64VarP(&bufferSizeArg, "buffer-size", "b", bufferSizeArg,
		"bidirectional buffer size in bytes")

	_ = cmd.MarkPersistentFlagRequired("downstream")
	_ = cmd.MarkPersistentFlagRequired("upstream")
	_ = cmd.MarkPersistentFlagRequired("buffer-size")

	return cmd
}
