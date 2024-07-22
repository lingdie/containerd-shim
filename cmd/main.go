package main

import (
	"cri-shim/pkg/shim"
	"cri-shim/pkg/types"

	shimapi "github.com/containerd/containerd/runtime/v2/shim"
)

func shimConfig(config *shimapi.Config) {
	config.NoReaper = true
	config.NoSubreaper = true
}

func main() {
	shimapi.Run(types.DefaultShimName, shim.New, shimConfig)
}
