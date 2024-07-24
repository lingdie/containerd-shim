package main

import (
	imageutil "cri-shim/pkg/image"
	"flag"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"cri-shim/pkg/server"
)

var criSocket, shimSocket string
var debug bool

func main() {
	flag.StringVar(&criSocket, "cri-socket", "unix:///var/run/containerd/containerd.sock", "CRI socket path")
	flag.StringVar(&shimSocket, "shim-socket", "/var/run/sealos/containerd-shim.sock", "CRI shim socket path")
	flag.StringVar(&imageutil.DefaultUserName, "username", "", "Image Hub username")
	flag.StringVar(&imageutil.DefaultPassword, "password", "", "Image Hub password")
	flag.BoolVar(&debug, "debug", false, "enable debug logging")
	flag.Parse()

	if imageutil.DefaultUserName == "" {
		slog.Error("failed to get username", nil)
		return
	}

	s, err := server.New(server.Options{
		Timeout:    time.Minute * 5,
		ShimSocket: shimSocket,
		CRISocket:  criSocket,
	})
	if debug {
		slog.SetLogLoggerLevel(slog.LevelDebug)
	}
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
