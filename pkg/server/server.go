package server

import (
	"context"
	"fmt"
	"net"
	"os"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	netutil "cri-shim/pkg/net"
	runtimeapi "k8s.io/cri-api/pkg/apis/runtime/v1"
)

type Options struct {
	Timeout    time.Duration
	ShimSocket string
	CRISocket  string
	// User is the user ID for our gRPC socket.
	User int
	// Group is the group ID for our gRPC socket.
	Group int
	// Mode is the permission mode bits for our gRPC socket.
	Mode os.FileMode
}

type Server struct {
	client   runtimeapi.RuntimeServiceClient
	server   *grpc.Server
	listener net.Listener
	options  Options
}

func New(options Options) (*Server, error) {
	listener, err := net.Listen("unix", options.ShimSocket)
	if err != nil {
		return nil, err
	}
	server := grpc.NewServer()
	return &Server{
		server:   server,
		listener: listener,
		options:  options,
	}, nil
}

func (s *Server) Start() error {
	go func() {
		_ = s.server.Serve(s.listener)
	}()
	conn, err := grpc.NewClient(s.options.CRISocket, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return err
	}
	s.client = runtimeapi.NewRuntimeServiceClient(conn)
	runtimeapi.RegisterRuntimeServiceServer(s.server, s)
	return netutil.WaitForServer(s.options.ShimSocket, time.Second)
}

func (s *Server) Stop() {
	s.server.Stop()
	s.listener.Close()
}

func (s *Server) Version(ctx context.Context, request *runtimeapi.VersionRequest) (*runtimeapi.VersionResponse, error) {
	return s.client.Version(ctx, request)
}

func (s *Server) RunPodSandbox(ctx context.Context, request *runtimeapi.RunPodSandboxRequest) (*runtimeapi.RunPodSandboxResponse, error) {
	return s.client.RunPodSandbox(ctx, request)
}

func (s *Server) StopPodSandbox(ctx context.Context, request *runtimeapi.StopPodSandboxRequest) (*runtimeapi.StopPodSandboxResponse, error) {
	return s.client.StopPodSandbox(ctx, request)
}

func (s *Server) RemovePodSandbox(ctx context.Context, request *runtimeapi.RemovePodSandboxRequest) (*runtimeapi.RemovePodSandboxResponse, error) {
	return s.client.RemovePodSandbox(ctx, request)
}

func (s *Server) PodSandboxStatus(ctx context.Context, request *runtimeapi.PodSandboxStatusRequest) (*runtimeapi.PodSandboxStatusResponse, error) {
	return s.client.PodSandboxStatus(ctx, request)
}

func (s *Server) ListPodSandbox(ctx context.Context, request *runtimeapi.ListPodSandboxRequest) (*runtimeapi.ListPodSandboxResponse, error) {
	return s.client.ListPodSandbox(ctx, request)
}

func (s *Server) CreateContainer(ctx context.Context, request *runtimeapi.CreateContainerRequest) (*runtimeapi.CreateContainerResponse, error) {
	return s.client.CreateContainer(ctx, request)
}

func (s *Server) StartContainer(ctx context.Context, request *runtimeapi.StartContainerRequest) (*runtimeapi.StartContainerResponse, error) {
	return s.client.StartContainer(ctx, request)
}

func (s *Server) StopContainer(ctx context.Context, request *runtimeapi.StopContainerRequest) (*runtimeapi.StopContainerResponse, error) {
	// todo check container env and create commit
	return s.client.StopContainer(ctx, request)
}

func (s *Server) RemoveContainer(ctx context.Context, request *runtimeapi.RemoveContainerRequest) (*runtimeapi.RemoveContainerResponse, error) {
	return s.client.RemoveContainer(ctx, request)
}

func (s *Server) ListContainers(ctx context.Context, request *runtimeapi.ListContainersRequest) (*runtimeapi.ListContainersResponse, error) {
	return s.client.ListContainers(ctx, request)
}

func (s *Server) ContainerStatus(ctx context.Context, request *runtimeapi.ContainerStatusRequest) (*runtimeapi.ContainerStatusResponse, error) {
	return s.client.ContainerStatus(ctx, request)
}

func (s *Server) UpdateContainerResources(ctx context.Context, request *runtimeapi.UpdateContainerResourcesRequest) (*runtimeapi.UpdateContainerResourcesResponse, error) {
	return s.client.UpdateContainerResources(ctx, request)
}

func (s *Server) ReopenContainerLog(ctx context.Context, request *runtimeapi.ReopenContainerLogRequest) (*runtimeapi.ReopenContainerLogResponse, error) {
	return s.client.ReopenContainerLog(ctx, request)
}

func (s *Server) ExecSync(ctx context.Context, request *runtimeapi.ExecSyncRequest) (*runtimeapi.ExecSyncResponse, error) {
	return s.client.ExecSync(ctx, request)
}

func (s *Server) Exec(ctx context.Context, request *runtimeapi.ExecRequest) (*runtimeapi.ExecResponse, error) {
	return s.client.Exec(ctx, request)
}

func (s *Server) Attach(ctx context.Context, request *runtimeapi.AttachRequest) (*runtimeapi.AttachResponse, error) {
	return s.client.Attach(ctx, request)
}

func (s *Server) PortForward(ctx context.Context, request *runtimeapi.PortForwardRequest) (*runtimeapi.PortForwardResponse, error) {
	return s.client.PortForward(ctx, request)
}

func (s *Server) ContainerStats(ctx context.Context, request *runtimeapi.ContainerStatsRequest) (*runtimeapi.ContainerStatsResponse, error) {
	return s.client.ContainerStats(ctx, request)
}

func (s *Server) ListContainerStats(ctx context.Context, request *runtimeapi.ListContainerStatsRequest) (*runtimeapi.ListContainerStatsResponse, error) {
	return s.client.ListContainerStats(ctx, request)
}

func (s *Server) PodSandboxStats(ctx context.Context, request *runtimeapi.PodSandboxStatsRequest) (*runtimeapi.PodSandboxStatsResponse, error) {
	return s.client.PodSandboxStats(ctx, request)
}

func (s *Server) ListPodSandboxStats(ctx context.Context, request *runtimeapi.ListPodSandboxStatsRequest) (*runtimeapi.ListPodSandboxStatsResponse, error) {
	return s.client.ListPodSandboxStats(ctx, request)
}

func (s *Server) UpdateRuntimeConfig(ctx context.Context, request *runtimeapi.UpdateRuntimeConfigRequest) (*runtimeapi.UpdateRuntimeConfigResponse, error) {
	return s.client.UpdateRuntimeConfig(ctx, request)
}

func (s *Server) Status(ctx context.Context, request *runtimeapi.StatusRequest) (*runtimeapi.StatusResponse, error) {
	return s.client.Status(ctx, request)
}

func (s *Server) CheckpointContainer(ctx context.Context, request *runtimeapi.CheckpointContainerRequest) (*runtimeapi.CheckpointContainerResponse, error) {
	return s.client.CheckpointContainer(ctx, request)
}

func (s *Server) GetContainerEvents(request *runtimeapi.GetEventsRequest, server runtimeapi.RuntimeService_GetContainerEventsServer) error {
	client, err := s.client.GetContainerEvents(context.Background(), request)
	if err != nil {
		return err
	}
	if res, err := client.Recv(); err != nil {
		return err
	} else {
		return server.Send(res)
	}
}

func (s *Server) ListMetricDescriptors(ctx context.Context, request *runtimeapi.ListMetricDescriptorsRequest) (*runtimeapi.ListMetricDescriptorsResponse, error) {
	return s.client.ListMetricDescriptors(ctx, request)
}

func (s *Server) ListPodSandboxMetrics(ctx context.Context, request *runtimeapi.ListPodSandboxMetricsRequest) (*runtimeapi.ListPodSandboxMetricsResponse, error) {
	return s.client.ListPodSandboxMetrics(ctx, request)
}

func (s *Server) RuntimeConfig(ctx context.Context, request *runtimeapi.RuntimeConfigRequest) (*runtimeapi.RuntimeConfigResponse, error) {
	return s.client.RuntimeConfig(ctx, request)
}

func serverError(format string, args ...interface{}) error {
	return fmt.Errorf("cri-shim/server: "+format, args...)
}
