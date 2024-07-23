package main

import (
	"flag"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"cri-shim/pkg/server"
)

var criSocket, shimSocket, timeout string

func main() {
	flag.StringVar(&criSocket, "cri-socket", "unix:///var/run/containerd/containerd.sock", "CRI socket path")
	flag.StringVar(&shimSocket, "shim-socket", "/var/run/sealos/containerd-shim.sock", "CRI shim socket path")
	flag.Parse()
	s, err := server.New(server.Options{
		Timeout:    time.Minute * 5,
		ShimSocket: shimSocket,
		CRISocket:  criSocket,
	})
	if err != nil {
		slog.Error("failed to create server", err)
		return
	}
	err = s.Start()
	if err != nil {
		slog.Error("failed to start server", err)
		return
	}
	slog.Info("server started")

	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh, syscall.SIGINT, syscall.SIGTERM)

	stopCh := make(chan struct{}, 1)
	select {
	case <-signalCh:
		close(stopCh)
	case <-stopCh:
	}
	_ = os.Remove(shimSocket)
	slog.Info("shutting down the image_shim")
}
