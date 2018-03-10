package main

import (
	"context"
	"fmt"
	"io"
	"os"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/client"
)

func main() {
	ctx := context.Background()

	client, err := client.NewEnvClient()
	if err != nil {
		panic(err)
	}

	config := container.Config{
		Image: "library/alpine",
		Cmd:   []string{"ping", "-c", "3", "8.8.8.8"},
	}

	res, err := client.ContainerCreate(ctx, &config, nil, nil, "")
	if err != nil {
		panic(err)
	}

	containerID := res.ID
	fmt.Println("Container ID is ", res.ID)

	// Start the container

	fmt.Println("Starting container ", containerID)

	if err := client.ContainerStart(ctx, containerID, types.ContainerStartOptions{}); err != nil {
		panic(err)
	}

	// Wait for the container to finish

	fmt.Println("Waiting for container to finish")

	resultc, errorc := client.ContainerWait(ctx, containerID, container.WaitConditionRemoved)

	select {
	case err := <-errorc:
		fmt.Println("Got an error: ", err)
	case res := <-resultc:
		fmt.Println("Got a result: ", res.StatusCode)
	}

	//

	fmt.Println("Container logs:")

	containerLogsOptions := types.ContainerLogsOptions{
		ShowStdout: true,
		ShowStderr: true,
		Timestamps: true,
	}

	out, err := client.ContainerLogs(ctx, containerID, containerLogsOptions)
	if err != nil {
		panic(err)
	}

	if _, err := io.Copy(os.Stdout, out); err != nil {
		panic(err)
	}

	// Prune containers

	fmt.Println("Pruning containers")

	if _, err := client.ContainersPrune(ctx, filters.Args{}); err != nil {
		panic(err)
	}
}
