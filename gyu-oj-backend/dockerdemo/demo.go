package main

import (
	"context"
	"fmt"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/image"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/client"
	ocispec "github.com/opencontainers/image-spec/specs-go/v1"
	"io"
	"log"
	"os"
	"sync"
	"time"
)

/*
1，拉取镜像
2，创建容器
3，查看容器状态
4，启动容器
5，查看日志
6，删除容器
7，删除镜像
*/

var wg sync.WaitGroup

func DockerDemo() {
	// 1, 创建一个与 Docker 通信的客户端
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		log.Fatal(err)
	}

	defer cli.Close()

	image := "nginx:latest"

	// 拉取镜像的写法
	// 可能拉取镜像操作耗时较长，所以开一个 goroutine 并监听到结束拉取成功为止
	//wg.Add(1)
	//go func() {
	//	err = PullImage(context.Background(), cli, "nginx:latest")
	//	if err != nil {
	//		log.Fatal(err)
	//	}
	//	wg.Done()
	//}()
	//wg.Wait()

	// 列出所有容器
	//containers, err := cli.ContainerList(context.Background(), container.ListOptions{All: true})
	//if err != nil {
	//	log.Fatal(err)
	//}
	//// 输出容器信息
	//fmt.Println("containerId containerImage")
	//for _, container := range containers {
	//	fmt.Printf("%s %s\n", container.ID[:10], container.Image)
	//}

	// 2, 创建容器
	args := CreateContainerArgs{
		config: &container.Config{
			Cmd:   []string{"echo", "Hello Docker"},
			Image: image,
		},
		hostConfig:       nil,
		networkingConfig: nil,
		platform:         nil,
		containerName:    "",
	}
	containerId, err := CreateContainer(context.Background(), cli, args)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("创建的容器 ID: ", containerId)

	// 3,启动容器
	err = cli.ContainerStart(context.Background(), containerId, container.StartOptions{})
	if err != nil {
		log.Fatal(err)
	}

	// 4,查看容器日志
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	go func(ctx context.Context, containerId string) {
		reader, err := cli.ContainerLogs(ctx, containerId, container.LogsOptions{
			ShowStdout: true,
			ShowStderr: true,
		})
		if err != nil {
			log.Fatal(err)
		}

		_, err = io.Copy(os.Stdout, reader)
		if err != nil && err != io.EOF {
			log.Fatal(err)
		}
	}(ctx, containerId)

	time.Sleep(7 * time.Second)

	// 5, 删除容器
	err = cli.ContainerRemove(context.Background(), containerId, container.RemoveOptions{})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("删除容器成功")

}

// 拉取镜像

func PullImage(ctx context.Context, cli *client.Client, imageName string) error {
	pullImageReader, err := cli.ImagePull(ctx, imageName, image.PullOptions{})
	if err != nil {
		return err
	}
	_, err = io.Copy(os.Stdout, pullImageReader)
	if err != nil && err != io.EOF {
		log.Fatal(err)
	}
	return nil
}

// 创建容器

func CreateContainer(ctx context.Context, cli *client.Client, args CreateContainerArgs) (containerId string, error error) {
	// 创建容器
	crateResp, err := cli.ContainerCreate(ctx, args.config, args.hostConfig, args.networkingConfig, args.platform, args.containerName)
	if err != nil {
		return "", err
	}
	// 查看当前所有容器状态
	list, err := cli.ContainerList(ctx, container.ListOptions{All: true})
	if err != nil {
		return "", err
	}
	for _, c := range list {
		fmt.Println("容器 ID: ", c.ID[:10])
	}

	return crateResp.ID, nil
}

type CreateContainerArgs struct {
	config           *container.Config
	hostConfig       *container.HostConfig
	networkingConfig *network.NetworkingConfig
	platform         *ocispec.Platform
	containerName    string
}
