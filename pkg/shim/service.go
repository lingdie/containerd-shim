package shim

import (
	"context"
	"fmt"
	"io"
	"os"

	taskAPI "github.com/containerd/containerd/api/runtime/task/v2"
	"github.com/containerd/containerd/errdefs"
	"github.com/containerd/containerd/namespaces"
	"github.com/containerd/containerd/runtime/v2/shim"

	"github.com/sirupsen/logrus"
	"google.golang.org/protobuf/types/known/emptypb"
)

// shimLog is logger for shim package
var shimLog = logrus.WithFields(logrus.Fields{
	"source": "containerd-shim-v2-sealos",
	"name":   "containerd-shim-v2",
})

func New(ctx context.Context, id string, publisher shim.Publisher, shutdown func()) (shim.Shim, error) {
	shimLog = shimLog.WithFields(logrus.Fields{
		"sandbox": id,
		"pid":     os.Getpid(),
	})
	// Discard the log before shim init its log output. Otherwise
	// it will output into stdio, from which containerd would like
	// to get the shim's socket address.
	logrus.SetOutput(io.Discard)
	opts := ctx.Value(shim.OptsKey{}).(shim.Opts)
	if !opts.Debug {
		logrus.SetLevel(logrus.WarnLevel)
	}

	ns, found := namespaces.Namespace(ctx)
	if !found {
		return nil, fmt.Errorf("shim namespace cannot be empty")
	}

	s := &service{
		id:        id,
		ctx:       ctx,
		cancel:    shutdown,
		namespace: ns,
	}

	return s, nil
}

// service is the shim implementation of a remote shim over GRPC
type service struct {
	ctx       context.Context
	cancel    func()
	id        string
	namespace string
}

// StartShim is a binary call that executes a new shim returning the address
func (s *service) StartShim(ctx context.Context, opts shim.StartOpts) (string, error) {
	return "", nil
}

// Cleanup is a binary call that cleans up any resources used by the shim when the service crashes
func (s *service) Cleanup(ctx context.Context) (*taskAPI.DeleteResponse, error) {
	return nil, errdefs.ErrNotImplemented
}

// Create a new container
func (s *service) Create(ctx context.Context, r *taskAPI.CreateTaskRequest) (_ *taskAPI.CreateTaskResponse, err error) {
	return nil, errdefs.ErrNotImplemented
}

// Start the primary user process inside the container
func (s *service) Start(ctx context.Context, r *taskAPI.StartRequest) (*taskAPI.StartResponse, error) {
	return nil, errdefs.ErrNotImplemented
}

// Delete a process or container
func (s *service) Delete(ctx context.Context, r *taskAPI.DeleteRequest) (*taskAPI.DeleteResponse, error) {
	return nil, errdefs.ErrNotImplemented
}

// Exec an additional process inside the container
func (s *service) Exec(ctx context.Context, r *taskAPI.ExecProcessRequest) (*emptypb.Empty, error) {
	return nil, errdefs.ErrNotImplemented
}

// ResizePty of a process
func (s *service) ResizePty(ctx context.Context, r *taskAPI.ResizePtyRequest) (*emptypb.Empty, error) {
	return nil, errdefs.ErrNotImplemented
}

// State returns runtime state of a process
func (s *service) State(ctx context.Context, r *taskAPI.StateRequest) (*taskAPI.StateResponse, error) {
	return nil, errdefs.ErrNotImplemented
}

// Pause the container
func (s *service) Pause(ctx context.Context, r *taskAPI.PauseRequest) (*emptypb.Empty, error) {
	return nil, errdefs.ErrNotImplemented
}

// Resume the container
func (s *service) Resume(ctx context.Context, r *taskAPI.ResumeRequest) (*emptypb.Empty, error) {
	return nil, errdefs.ErrNotImplemented
}

// Kill a process
func (s *service) Kill(ctx context.Context, r *taskAPI.KillRequest) (*emptypb.Empty, error) {
	return nil, errdefs.ErrNotImplemented
}

// Pids returns all pids inside the container
func (s *service) Pids(ctx context.Context, r *taskAPI.PidsRequest) (*taskAPI.PidsResponse, error) {
	return nil, errdefs.ErrNotImplemented
}

// CloseIO of a process
func (s *service) CloseIO(ctx context.Context, r *taskAPI.CloseIORequest) (*emptypb.Empty, error) {
	return nil, errdefs.ErrNotImplemented
}

// Checkpoint the container
func (s *service) Checkpoint(ctx context.Context, r *taskAPI.CheckpointTaskRequest) (*emptypb.Empty, error) {
	return nil, errdefs.ErrNotImplemented
}

// Connect returns shim information of the underlying service
func (s *service) Connect(ctx context.Context, r *taskAPI.ConnectRequest) (*taskAPI.ConnectResponse, error) {
	return nil, errdefs.ErrNotImplemented
}

// Shutdown is called after the underlying resources of the shim are cleaned up and the service can be stopped
func (s *service) Shutdown(ctx context.Context, r *taskAPI.ShutdownRequest) (*emptypb.Empty, error) {
	os.Exit(0)
	return &emptypb.Empty{}, nil
}

// Stats returns container level system stats for a container and its processes
func (s *service) Stats(ctx context.Context, r *taskAPI.StatsRequest) (*taskAPI.StatsResponse, error) {
	return nil, errdefs.ErrNotImplemented
}

// Update the live container
func (s *service) Update(ctx context.Context, r *taskAPI.UpdateTaskRequest) (*emptypb.Empty, error) {
	return nil, errdefs.ErrNotImplemented
}

// Wait for a process to exit
func (s *service) Wait(ctx context.Context, r *taskAPI.WaitRequest) (*taskAPI.WaitResponse, error) {
	return nil, errdefs.ErrNotImplemented
}
