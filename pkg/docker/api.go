// Copyright (c) 2022 Institute of Software, Chinese Academy of Sciences (ISCAS)
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package docker

import (
	"context"
	"errors"
	"fmt"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/client"
	"io"
	"os"
	"strings"
)

var c *ClientAPI

type ClientAPI struct {
	DockerCli *client.Client
	Ctx       context.Context
}

type ClientAPIer interface {
	Stop(containerId string) error
	PullImage(imageUrl string) error
	Create(config *container.Config, hostConfig *container.HostConfig, networkingConfig *network.NetworkingConfig, containerName string) (string, error)
	Start(containerId string, startOpts types.ContainerStartOptions) error
	EnsureImageExist(image string) (bool, error)
	CreateNetwork(containerName string) (string, error)
	ConnectNetwork(networkId string, containerId string) error
	RemoveImage(imageId string, removeOpt types.ImageRemoveOptions) error
	Status(containerId string) (*types.Container, error)
	Inspect(containerId string) (*types.ContainerJSON, error)
	RemoveContainer(containerId string, removeVolumes bool) error
	RemoveNetwork(networkId string) error
	Exec(containerId string, cmd []string) error
	EnsureNetworkExist(networkName string) (string, error)
	FindContainer(containerName string) (string, error)
	IsImageUsed(imageID string) (bool, error)
}

func GetContainerAPIer() ClientAPIer {
	return c
}

func NewContainerMgr() ClientAPIer {
	return &ClientAPI{
		DockerCli: DockerClient(),
		Ctx:       context.Background(),
	}
}

func (c *ClientAPI) Create(config *container.Config, hostConfig *container.HostConfig, networkConfig *network.NetworkingConfig, containerName string) (string, error) {
	//var networkConfig *network.NetworkingConfig

	rspBody, err := c.DockerCli.ContainerCreate(c.Ctx, config, hostConfig, networkConfig, nil, containerName)
	if err != nil {
		return "", err
	}
	return rspBody.ID, nil
}

func (c *ClientAPI) Start(containerId string, startOpts types.ContainerStartOptions) error {
	err := c.DockerCli.ContainerStart(c.Ctx, containerId, startOpts)
	if err != nil {
		return err
	}
	return nil
}

func (c *ClientAPI) Status(containerId string) (*types.Container, error) {
	var args = filters.NewArgs()
	args.Add("id", containerId)
	containerInfo, err := c.DockerCli.ContainerList(c.Ctx, types.ContainerListOptions{
		All:     true,
		Filters: args})

	if err != nil {
		return nil, err
	}
	if len(containerInfo) > 0 {
		return &containerInfo[0], nil
	}
	return nil, errors.New("not found container")
}

func (c *ClientAPI) Inspect(containerId string) (*types.ContainerJSON, error) {
	containerInpsect, err := c.DockerCli.ContainerInspect(c.Ctx, containerId)
	if err != nil {
		return nil, err
	}
	return &containerInpsect, nil
}

func (c *ClientAPI) Stop(containerId string) error {
	duration := 10

	err := c.DockerCli.ContainerStop(c.Ctx, containerId, container.StopOptions{Timeout: &duration})

	if err != nil {
		return err
	}
	return nil
}

func (c *ClientAPI) RemoveContainer(containerId string, removeVolumes bool) error {
	//var duration = 10 * time.Second

	err := c.DockerCli.ContainerRemove(c.Ctx, containerId, types.ContainerRemoveOptions{RemoveVolumes: removeVolumes})

	if err != nil {
		return err
	}
	return nil
}

func (c *ClientAPI) PullImage(imageUrl string) error {
	//cm.DockerCli.ImageList(cm.Ctx,types.ImageListOptions{Filters: })
	events, err := c.DockerCli.ImagePull(c.Ctx, imageUrl, types.ImagePullOptions{})
	if err != nil {
		return err
	}
	defer events.Close()
	io.Copy(os.Stdout, events)
	return nil

}

func (c *ClientAPI) EnsureImageExist(image string) (bool, error) {
	var args = filters.NewArgs()
	if !strings.Contains(image, "@") {
		args.Add("reference", image)
	}
	images, err := c.DockerCli.ImageList(c.Ctx, types.ImageListOptions{Filters: args})
	if err != nil {
		return false, err
	}

	for _, _image := range images {
		//fmt.Println(_image.RepoDigests[0])
		if strings.Contains(image, "@") {
			if strings.Contains(_image.RepoDigests[0], image) {
				return true, nil
			}
		} else {
			if strings.Contains(_image.RepoTags[0], image) {
				return true, nil
			}
		}
	}
	return false, errors.New(fmt.Sprintf("can not find image %s", image))
}

func (c *ClientAPI) RemoveImage(imageId string, removeOpt types.ImageRemoveOptions) error {
	_, err := c.DockerCli.ImageRemove(c.Ctx, imageId, removeOpt)
	if err != nil {
		return err
	}
	return nil
}

func (c *ClientAPI) CreateNetwork(containerName string) (string, error) {
	createdNetwork, err := c.DockerCli.NetworkCreate(c.Ctx, containerName, types.NetworkCreate{})
	if err != nil {
		return "", err
	}
	return createdNetwork.ID, nil
}

func (c *ClientAPI) ConnectNetwork(networkId string, containerId string) error {
	var endpointSettings *network.EndpointSettings

	err := c.DockerCli.NetworkConnect(c.Ctx, networkId, containerId, endpointSettings)
	if err != nil {
		return err
	}
	return nil
}

func (c *ClientAPI) RemoveNetwork(networkId string) error {
	err := c.DockerCli.NetworkRemove(c.Ctx, networkId)
	if err != nil {
		return err
	}
	return nil
}

func (c *ClientAPI) EnsureNetworkExist(networkName string) (string, error) {
	var args = filters.NewArgs()
	args.Add("name", networkName)
	networkRes, err := c.DockerCli.NetworkList(c.Ctx, types.NetworkListOptions{Filters: args})
	if err != nil {
		return "", err
	}
	if len(networkRes) > 0 {
		for _, network := range networkRes {
			if network.Name == networkName {
				return network.ID, nil
			}
		}
		return "", nil
	} else {
		return "", nil
	}
}

func (c *ClientAPI) Exec(containerId string, cmd []string) error {
	execId, err := c.DockerCli.ContainerExecCreate(c.Ctx, containerId, types.ExecConfig{
		Cmd: cmd,
	})
	if err != nil {
		return err
	}
	err = c.DockerCli.ContainerExecStart(c.Ctx, execId.ID, types.ExecStartCheck{
		Detach: false,
		Tty:    false,
	})
	if err != nil {
		return err
	}
	return nil
}

func (c *ClientAPI) FindContainer(containerName string) (string, error) {
	var args = filters.NewArgs()
	args.Add("name", containerName)
	containerInfos, err := c.DockerCli.ContainerList(c.Ctx, types.ContainerListOptions{
		All:     true,
		Filters: args})

	if err != nil {
		return "", err
	}
	if len(containerInfos) > 0 {
		for i, ci := range containerInfos {
			if strings.TrimLeft(ci.Names[0], "/") == containerName {
				return containerInfos[i].ID, nil
			}
		}
	}
	return "", nil
}

func (c *ClientAPI) IsImageUsed(imageID string) (bool, error) {
	containerInfos, err := c.DockerCli.ContainerList(c.Ctx, types.ContainerListOptions{
		All: true})
	if err != nil {
		return false, err
	}
	for _, info := range containerInfos {
		if info.ImageID == imageID {
			return true, nil
		}
	}
	return false, nil
}
