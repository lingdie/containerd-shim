package image

import (
	"context"
	"github.com/containerd/nerdctl/v2/pkg/api/types"
	"github.com/containerd/nerdctl/v2/pkg/clientutil"
	"github.com/containerd/nerdctl/v2/pkg/cmd/container"
	"github.com/containerd/nerdctl/v2/pkg/cmd/image"
	"github.com/containerd/nerdctl/v2/pkg/cmd/login"
	"io"
)

// ImageInterface 定义了镜像操作接口
type ImageInterface interface {
	Push(args []string) error
	Commit(imageName, containerID string, pause bool) error
	Login(serverAddress, username, password string) error
}

// imageInterfaceImpl 实现了 ImageInterface 接口
type imageInterfaceImpl struct {
	GlobalOptions types.GlobalCommandOptions
	Stdout        io.Writer
}

// NewImageInterface 返回一个新的 ImageInterface 实现
// namespace: 命名空间
// address: 容器运行时的地址
// writer: 用于输出的 io.Writer
func NewImageInterface(namespace, address string, writer io.Writer) ImageInterface {
	return &imageInterfaceImpl{
		GlobalOptions: types.GlobalCommandOptions{
			Namespace: namespace,
			Address:   address,
		},
		Stdout: writer,
	}
}

// Commit 提交容器为镜像
// imageName: 镜像名称
// containerID: 容器 ID
// pause: 是否暂停容器后commit
func (impl *imageInterfaceImpl) Commit(imageName, containerID string, pause bool) error {
	options := types.ContainerCommitOptions{
		Stdout:   impl.Stdout,
		GOptions: impl.GlobalOptions,
		Pause:    pause,
	}

	client, ctx, cancel, err := clientutil.NewClient(context.Background(), options.GOptions.Namespace, options.GOptions.Address)
	if err != nil {
		return err
	}
	defer cancel()

	tmpName := imageName + "tmp"

	if err = container.Commit(ctx, client, tmpName, containerID, options); err != nil {
		return err
	}

	if err = impl.convert(tmpName, imageName); err != nil {
		return err
	}

	return impl.remove(tmpName, true, false)
}

// convert 将镜像转换为指定格式
// srcRawRef: 源镜像引用
// destRawRef: 目标镜像引用
func (impl *imageInterfaceImpl) convert(srcRawRef, destRawRef string) error {
	options := types.ImageConvertOptions{
		GOptions: impl.GlobalOptions,
		Oci:      true,
		Stdout:   impl.Stdout,
	}

	client, ctx, cancel, err := clientutil.NewClient(context.Background(), options.GOptions.Namespace, options.GOptions.Address)
	if err != nil {
		return err
	}
	defer cancel()

	return image.Convert(ctx, client, srcRawRef, destRawRef, options)
}

// remove 删除指定的镜像
// args: 镜像列表
// force: 是否强制删除
// async: 是否异步删除
func (impl *imageInterfaceImpl) remove(args string, force, async bool) error {
	options := types.ImageRemoveOptions{
		Stdout:   impl.Stdout,
		GOptions: impl.GlobalOptions,
		Force:    force,
		Async:    async,
	}

	client, ctx, cancel, err := clientutil.NewClient(context.Background(), options.GOptions.Namespace, options.GOptions.Address)
	if err != nil {
		return err
	}
	defer cancel()

	return image.Remove(ctx, client, []string{args}, options)
}

// Push 推送镜像到远程仓库
// args: 镜像列表
func (impl *imageInterfaceImpl) Push(args []string) error {
	options := types.ImagePushOptions{
		GOptions: impl.GlobalOptions,
		Stdout:   impl.Stdout,
	}

	client, ctx, cancel, err := clientutil.NewClient(context.Background(), options.GOptions.Namespace, options.GOptions.Address)
	if err != nil {
		return err
	}
	defer cancel()

	return image.Push(ctx, client, args[0], options)
}

// Login 登录到镜像仓库
// serverAddress: 仓库地址
// username: 用户名
// password: 密码
func (impl *imageInterfaceImpl) Login(serverAddress, username, password string) error {
	options := types.LoginCommandOptions{
		GOptions: impl.GlobalOptions,
		Username: username,
		Password: password,
	}
	if serverAddress != "" {
		options.ServerAddress = serverAddress
	}

	return login.Login(context.Background(), options, impl.Stdout)
}
