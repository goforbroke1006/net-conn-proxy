package cmd

import (
	"context"
	"fmt"
	"net"
	"os/signal"

	"github.com/spf13/cobra"
	"go.uber.org/zap"

	"github.com/goforbroke1006/net-conn-proxy/internal/component/tcp"
)

// NewTCPCmd represents the tcp command
func NewTCPCmd() *cobra.Command {
	var (
		downstreamAddrArg = "0.0.0.0:0"
		upstreamAddrArg   string
		bufferSizeArg     uint64 = 2048
	)

	cmd := &cobra.Command{
		Use: "tcp",
		Run: func(cmd *cobra.Command, args []string) {
			ctx, stop := signal.NotifyContext(context.Background())

			listener, err := net.Listen("tcp", downstreamAddrArg)
			if err != nil {
				zap.L().Fatal("downstream TCP", zap.Error(err))
			}
			defer func() { _ = listener.Close() }()

			zap.L().Info("TCP proxying",
				zap.String("up", upstreamAddrArg),
				zap.Uint64("buffer", bufferSizeArg))

			go func() {
				for {
					clientConn, err := listener.Accept()
					if err != nil {
						zap.L().Error("accept tcp conn", zap.Error(err))
						continue
					}
					zap.L().Info("accept tcp conn", zap.String("ip", clientConn.RemoteAddr().String()))
					go tcp.InitTCPProxy(clientConn, upstreamAddrArg, bufferSizeArg)
				}
			}()

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
