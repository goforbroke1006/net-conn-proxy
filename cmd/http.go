package cmd

import (
	"context"
	"fmt"
	"net"
	"os/signal"

	"github.com/spf13/cobra"
	"go.uber.org/zap"

	"github.com/goforbroke1006/net-conn-proxy/internal/component/http"
)

// NewHTTPCmd represents the http command
func NewHTTPCmd() *cobra.Command {
	var (
		downstreamAddrArg = "0.0.0.0:8080"
		originAddrArg     string
		bufferSizeArg     uint64 = 2048
	)

	cmd := &cobra.Command{
		Use: "http",
		Run: func(cmd *cobra.Command, args []string) {
			ctx, stop := signal.NotifyContext(context.Background())

			listener, err := net.Listen("tcp", downstreamAddrArg)
			if err != nil {
				zap.L().Fatal("downstream HTTP", zap.Error(err))
			}
			defer func() { _ = listener.Close() }()

			zap.L().Info("HTTP proxying",
				zap.String("origin", originAddrArg),
				zap.Uint64("buffer", bufferSizeArg))

			go func() {
				for {
					clientConn, err := listener.Accept()
					if err != nil {
						zap.L().Error("accept tcp conn", zap.Error(err))
						continue
					}
					zap.L().Info("accept tcp conn", zap.String("ip", clientConn.RemoteAddr().String()))
					go http.InitHTTPProxy(clientConn, downstreamAddrArg, originAddrArg, bufferSizeArg)
				}
			}()

			<-ctx.Done()
			defer stop()
			fmt.Println("bye...")
		},
	}

	cmd.PersistentFlags().StringVarP(&downstreamAddrArg, "downstream", "d", downstreamAddrArg,
		"downstream addr like 120.0.0.1:8080")
	cmd.PersistentFlags().StringVarP(&originAddrArg, "origin", "o", originAddrArg,
		"upstream addr like 8.8.8.8:80")
	cmd.PersistentFlags().Uint64VarP(&bufferSizeArg, "buffer-size", "b", bufferSizeArg,
		"bidirectional buffer size in bytes")

	_ = cmd.MarkPersistentFlagRequired("downstream")
	_ = cmd.MarkPersistentFlagRequired("origin")
	_ = cmd.MarkPersistentFlagRequired("buffer-size")

	return cmd
}
