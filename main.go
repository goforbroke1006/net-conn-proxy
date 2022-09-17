package main

import (
	"go.uber.org/zap"

	"github.com/goforbroke1006/net-conn-proxy/cmd"
)

func main() {
	logger, _ := zap.NewProduction()
	defer func() { _ = logger.Sync() }()
	zap.ReplaceGlobals(logger)

	cmd.Execute()
}
