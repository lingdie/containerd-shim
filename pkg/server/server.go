package server

import (
	"context"
	"cri-shim/pkg/container"
	imageutil "cri-shim/pkg/image"
	netutil "cri-shim/pkg/net"
	"cri-shim/pkg/types"
	"encoding/json"
	"github.com/containerd/containerd/namespaces"
	"log/slog"
	"net"
	"os"
	"strings"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/test/bufconn"
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
	client      runtimeapi.RuntimeServiceClient
	server      *grpc.Server
	listener    net.Listener
	options     Options
	bufListener *bufconn.Listener
	imageClient imageutil.ImageInterface
}

func New(options Options) (*Server, error) {
	listener, err := net.Listen("unix", options.ShimSocket)
	if err != nil {
		return nil, err
	}
	server := grpc.NewServer()

	imageClient, err := imageutil.NewImageInterface(imageutil.DefaultNamespace, options.CRISocket, os.Stdout)
	if err != nil {
		return nil, err
	}
	return &Server{
		server:      server,
		listener:    listener,
		options:     options,
		imageClient: imageClient,
	}, nil
}

func (s *Server) Start() error {
	conn, err := grpc.NewClient(s.options.CRISocket, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return err
	}
	s.client = runtimeapi.NewRuntimeServiceClient(conn)
	runtimeapi.RegisterRuntimeServiceServer(s.server, s)

	// do serve after client is created and registered
	go func() {
		_ = s.server.Serve(s.listener)
	}()
	return netutil.WaitForServer(s.options.ShimSocket, time.Second)
}

func (s *Server) Stop() {
	s.server.Stop()
	s.listener.Close()
	s.imageClient.Stop()
}

func (s *Server) Version(ctx context.Context, request *runtimeapi.VersionRequest) (*runtimeapi.VersionResponse, error) {
	slog.Info("Doing version request", "request", request)
	resp, err := s.client.Version(ctx, request)
	if err != nil {
		slog.Error("failed to get version", "error", err)
		return resp, err
	}
	slog.Debug("Got version response", "response", resp)
	return resp, err
}

func (s *Server) RunPodSandbox(ctx context.Context, request *runtimeapi.RunPodSandboxRequest) (*runtimeapi.RunPodSandboxResponse, error) {
	slog.Info("Doing run pod sandbox request", "request", request)
	return s.client.RunPodSandbox(ctx, request)
}

func (s *Server) StopPodSandbox(ctx context.Context, request *runtimeapi.StopPodSandboxRequest) (*runtimeapi.StopPodSandboxResponse, error) {
	slog.Info("Doing stop pod sandbox request", "request", request)
	return s.client.StopPodSandbox(ctx, request)
}

func (s *Server) RemovePodSandbox(ctx context.Context, request *runtimeapi.RemovePodSandboxRequest) (*runtimeapi.RemovePodSandboxResponse, error) {
	slog.Info("Doing remove pod sandbox request", "request", request)
	return s.client.RemovePodSandbox(ctx, request)
}

func (s *Server) PodSandboxStatus(ctx context.Context, request *runtimeapi.PodSandboxStatusRequest) (*runtimeapi.PodSandboxStatusResponse, error) {
	slog.Info("Doing pod sandbox status request", "request", request)
	return s.client.PodSandboxStatus(ctx, request)
}

func (s *Server) ListPodSandbox(ctx context.Context, request *runtimeapi.ListPodSandboxRequest) (*runtimeapi.ListPodSandboxResponse, error) {
	slog.Info("Doing list pod sandbox request", "request", request)
	return s.client.ListPodSandbox(ctx, request)
}

func (s *Server) CreateContainer(ctx context.Context, request *runtimeapi.CreateContainerRequest) (*runtimeapi.CreateContainerResponse, error) {
	slog.Info("Doing create container request", "request", request)
	return s.client.CreateContainer(ctx, request)
}

func (s *Server) StartContainer(ctx context.Context, request *runtimeapi.StartContainerRequest) (*runtimeapi.StartContainerResponse, error) {
	slog.Info("Doing start container request", "request", request)
	return s.client.StartContainer(ctx, request)
}

func (s *Server) StopContainer(ctx context.Context, request *runtimeapi.StopContainerRequest) (*runtimeapi.StopContainerResponse, error) {
	// todo check container env and create commit
	slog.Info("Doing stop container request", "request", request)

	return s.client.StopContainer(ctx, request)
}

func (s *Server) RemoveContainer(ctx context.Context, request *runtimeapi.RemoveContainerRequest) (*runtimeapi.RemoveContainerResponse, error) {
	slog.Debug("Doing remove container request", "request", request)

	statusReq := &runtimeapi.ContainerStatusRequest{
		ContainerId: request.ContainerId,
		Verbose:     true,
	}
	statusResp, err := s.client.ContainerStatus(ctx, statusReq)
	if err != nil {
		slog.Error("failed to get container status", "error", err)
		return nil, err
	}

	registry, imageRef, flag, err := ContainerEnv(statusResp)

	//commit image
	if flag && err == nil {
		// todo report failed to commit containers
		// skip commit if container is not running
		if statusResp.Status.State != runtimeapi.ContainerState_CONTAINER_RUNNING {
			// do something, should we remove container if we can't commit it?
		}

		ctx = namespaces.WithNamespace(ctx, imageutil.DefaultNamespace)

		if err = s.imageClient.Login(ctx, registry.LoginAddress, registry.UserName, registry.Password); err != nil {
			slog.Error("failed to login register", "error", err)
			return nil, err
		}

		if err = s.imageClient.Commit(ctx, registry.GetImageRef(imageRef), statusResp.Status.Id, false); err != nil {
			slog.Error("failed to commit container", "error", err)
			return nil, err
		}
	}

	return s.client.RemoveContainer(ctx, request)
}

func (s *Server) ListContainers(ctx context.Context, request *runtimeapi.ListContainersRequest) (*runtimeapi.ListContainersResponse, error) {
	slog.Info("Doing list containers request", "request", request)
	return s.client.ListContainers(ctx, request)
}

func (s *Server) ContainerStatus(ctx context.Context, request *runtimeapi.ContainerStatusRequest) (*runtimeapi.ContainerStatusResponse, error) {
	slog.Info("Doing container status request", "request", request)
	request.Verbose = true
	resp, err := s.client.ContainerStatus(ctx, request)
	if err != nil {
		slog.Error("failed to get container status", "error", err)
		return resp, err
	}
	slog.Debug("Got container status response", "response", resp)
	return resp, err
}

func (s *Server) UpdateContainerResources(ctx context.Context, request *runtimeapi.UpdateContainerResourcesRequest) (*runtimeapi.UpdateContainerResourcesResponse, error) {
	slog.Info("Doing update container resources request", "request", request)
	return s.client.UpdateContainerResources(ctx, request)
}

func (s *Server) ReopenContainerLog(ctx context.Context, request *runtimeapi.ReopenContainerLogRequest) (*runtimeapi.ReopenContainerLogResponse, error) {
	slog.Info("Doing reopen container log request", "request", request)
	return s.client.ReopenContainerLog(ctx, request)
}

func (s *Server) ExecSync(ctx context.Context, request *runtimeapi.ExecSyncRequest) (*runtimeapi.ExecSyncResponse, error) {
	slog.Info("Doing exec sync request", "request", request)
	return s.client.ExecSync(ctx, request)
}

func (s *Server) Exec(ctx context.Context, request *runtimeapi.ExecRequest) (*runtimeapi.ExecResponse, error) {
	slog.Info("Doing exec request", "request", request)
	return s.client.Exec(ctx, request)
}

func (s *Server) Attach(ctx context.Context, request *runtimeapi.AttachRequest) (*runtimeapi.AttachResponse, error) {
	slog.Info("Doing attach request", "request", request)
	return s.client.Attach(ctx, request)
}

func (s *Server) PortForward(ctx context.Context, request *runtimeapi.PortForwardRequest) (*runtimeapi.PortForwardResponse, error) {
	slog.Info("Doing port forward request", "request", request)
	return s.client.PortForward(ctx, request)
}

func (s *Server) ContainerStats(ctx context.Context, request *runtimeapi.ContainerStatsRequest) (*runtimeapi.ContainerStatsResponse, error) {
	slog.Info("Doing container stats request", "request", request)
	return s.client.ContainerStats(ctx, request)
}

func (s *Server) ListContainerStats(ctx context.Context, request *runtimeapi.ListContainerStatsRequest) (*runtimeapi.ListContainerStatsResponse, error) {
	slog.Info("Doing list container stats request", "request", request)
	return s.client.ListContainerStats(ctx, request)
}

func (s *Server) PodSandboxStats(ctx context.Context, request *runtimeapi.PodSandboxStatsRequest) (*runtimeapi.PodSandboxStatsResponse, error) {
	slog.Info("Doing pod sandbox stats request", "request", request)
	return s.client.PodSandboxStats(ctx, request)
}

func (s *Server) ListPodSandboxStats(ctx context.Context, request *runtimeapi.ListPodSandboxStatsRequest) (*runtimeapi.ListPodSandboxStatsResponse, error) {
	slog.Info("Doing list pod sandbox stats request", "request", request)
	return s.client.ListPodSandboxStats(ctx, request)
}

func (s *Server) UpdateRuntimeConfig(ctx context.Context, request *runtimeapi.UpdateRuntimeConfigRequest) (*runtimeapi.UpdateRuntimeConfigResponse, error) {
	slog.Info("Doing update runtime config request", "request", request)
	return s.client.UpdateRuntimeConfig(ctx, request)
}

func (s *Server) Status(ctx context.Context, request *runtimeapi.StatusRequest) (*runtimeapi.StatusResponse, error) {
	slog.Info("Doing status request", "request", request)
	return s.client.Status(ctx, request)
}

func (s *Server) CheckpointContainer(ctx context.Context, request *runtimeapi.CheckpointContainerRequest) (*runtimeapi.CheckpointContainerResponse, error) {
	slog.Info("Doing checkpoint container request", "request", request)
	return s.client.CheckpointContainer(ctx, request)
}

func (s *Server) GetContainerEvents(request *runtimeapi.GetEventsRequest, server runtimeapi.RuntimeService_GetContainerEventsServer) error {
	slog.Info("Doing get container events request", "request", request)
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
	slog.Info("Doing list metric descriptors request", "request", request)
	return s.client.ListMetricDescriptors(ctx, request)
}

func (s *Server) ListPodSandboxMetrics(ctx context.Context, request *runtimeapi.ListPodSandboxMetricsRequest) (*runtimeapi.ListPodSandboxMetricsResponse, error) {
	slog.Info("Doing list pod sandbox metrics request", "request", request)
	return s.client.ListPodSandboxMetrics(ctx, request)
}

func (s *Server) RuntimeConfig(ctx context.Context, request *runtimeapi.RuntimeConfigRequest) (*runtimeapi.RuntimeConfigResponse, error) {
	slog.Info("Doing runtime config request", "request", request)
	return s.client.RuntimeConfig(ctx, request)
}

func ContainerEnv(resp *runtimeapi.ContainerStatusResponse) (*imageutil.Registry, string, bool, error) {
	info := &container.Info{}
	if err := json.Unmarshal([]byte(resp.Info["info"]), info); err != nil {
		slog.Error("failed to unmarshal container info", "error", err)
		return nil, "", false, err
	}
	slog.Debug("Got container info env", "info env", info.Config.Envs)
	var registryName, userName, password, imageName, repo string
	flag := false

	envMap := map[string]*string{
		types.ImageRegistryAddressOnEnv:    &registryName,
		types.ImageRegistryUserNameOnEnv:   &userName,
		types.ImageRegistryPasswordOnEnv:   &password,
		types.ImageNameOnEnv:               &imageName,
		types.ImageRegistryRepositoryOnEnv: &repo,
	}

	for _, env := range info.Config.Envs {
		if env.Key == types.ContainerCommitOnStopEnvFlag {
			if env.Value != types.ContainerCommitOnStopEnvEnableValue {
				return nil, "", false, nil
			}
			flag = true
			continue
		}

		if target, exists := envMap[env.Key]; exists {
			*target = env.Value
		}
	}

	if imageName == "" {
		imageName = resp.Status.Image.Image
		parts := strings.Split(imageName, "/")
		if len(parts) > 1 {
			imageName = strings.Join(parts[len(parts)-1:], "/")
		}
	}

	return imageutil.NewRegistry(registryName, repo, userName, password), imageName, flag, nil
}
