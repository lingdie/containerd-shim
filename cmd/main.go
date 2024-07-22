package main

import (
	"flag"
	"log/slog"
	"time"

	"cri-shim/pkg/server"
)

var criSocket, shimSocket, timeout string

func main() {
	flag.StringVar(&criSocket, "cri-socket", "/var/run/containerd/containerd.sock", "CRI socket path")
	flag.StringVar(&shimSocket, "shim-socket", "/var/run/sealos/cri-shim.sock", "CRI shim socket path")
	flag.Parse()
	s, err := server.New(server.Options{
		Timeout:    time.Minute * 5,
		ShimSocket: shimSocket,
		CRISocket:  criSocket,
	})
	if err != nil {
		slog.Error("failed to create server: %v", err)
		return
	}
	err = s.Start()
	if err != nil {
		slog.Error("failed to start server: %v", err)
		return
	}
}
