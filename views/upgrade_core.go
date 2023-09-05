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

package views

import (
	"errors"
	"eulixspace-upgrade/config"
	"eulixspace-upgrade/models"
	"eulixspace-upgrade/pkg/docker"
	"eulixspace-upgrade/utils"
	"eulixspace-upgrade/utils/db"
	"eulixspace-upgrade/utils/version"
	"fmt"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/go-connections/nat"
	"github.com/dungeonsnd/gocom/log4go"
	"os"
	"strconv"
	"strings"
	"time"
)

const agentRpmName = "eulixspace-agent"
const OsType = "aarch64"

type NativeAgent struct {
}

type ContainerAgent struct {
	DockerAPI docker.ClientAPIer
}

type Agent interface {
	Upgrade(versionId string) error
}

func (na *NativeAgent) Upgrade(versionId string) error {
	rpmPath, err := GetAgentRpmPath(versionId)
	if err != nil {
		log4go.E("Failed to get Agent rpmPath: %s", err)
		db.MarkTaskInstallErr(versionId)
		return err
	}
	log4go.I("install system-agent after 5 seconds")
	time.Sleep(5 * time.Second)
	err = dnfInstall(rpmPath)
	if err != nil {
		log4go.E("Failed to Install Agent %s", err)
		db.MarkTaskInstallErr(versionId)
		return err
	}
	return nil
}

func (ca *ContainerAgent) Upgrade(versionId string, dataDir string) error {
	// 容器环境下对all-in-one 容器进行升级
	// 停all-in-one 容器
	var imageUrl string
	cid, err := ca.DockerAPI.FindContainer(config.Config.Container.Name)
	if err != nil {
		return err
	}
	containerInfo, err := ca.DockerAPI.Status(cid)
	if err != nil {
		return err
	}
	// stop container
	log4go.I("stop aospace-all-in-one container")
	err = ca.DockerAPI.Stop(cid)
	if err != nil {
		return err
	}
	// remove container
	log4go.I("remove aospace-all-in-one container")
	err = ca.DockerAPI.RemoveContainer(cid, false)
	if err != nil {
		return err
	}
	// remove image
	log4go.I("remove aospace-all-in-one image")
	err = ca.DockerAPI.RemoveImage(containerInfo.ImageID, types.ImageRemoveOptions{})
	if err != nil {
		return err
	}
	// run new aospace-all-in-one
	var volumes []string
	var ports []models.Port
	var envs []models.Environment
	containerName := config.Config.Container.Name

	imageUrl = config.Config.Container.ImageOpenSource + ":" + versionId
	volumes = append(volumes, dataDir)
	volumes = append(volumes, "/var/run/docker.sock")
	webPort := models.Port{
		Number:   config.Config.Container.WebPort,
		Type:     "intranet",
		Usage:    "",
		Protocol: "tcp",
	}
	ApiPort := models.Port{
		Number:   config.Config.Container.ApiPort,
		Type:     "internal",
		Usage:    "",
		Protocol: "tcp",
	}
	ports = append(ports, ApiPort)
	ports = append(ports, webPort)

	dataDirEnv := models.Environment{
		Key:   "AOSPACE_DATADIR",
		Value: dataDir,
	}
	if !strings.Contains(utils.OSKernel(), "microsoft-standard-WSL") {
		networkModeEnv := models.Environment{
			Key:   "RUN_NETWORK_MODE",
			Value: "host",
		}
		envs = append(envs, networkModeEnv)
	}
	envs = append(envs, dataDirEnv)

	err = ca.RunContainer(dataDir, imageUrl, containerName, volumes, ports, envs)
	if err != nil {
		log4go.E("run container err:%v", err)
		return err
	}

	return nil
}

func (ca *ContainerAgent) RunContainer(dataDir string, imageUrl string, containerName string, volumes []string, ports []models.Port, envs []models.Environment) error {
	var (
		containerId  string
		environments []string
		portMap      = make(nat.PortMap)
		exposePort   = make(nat.PortSet)
		portsRsp     = make(map[int]int)
		hostIp       string
		err          error
	)
	networkConfig := &network.NetworkingConfig{}
	imageUrlTrim := strings.TrimSpace(imageUrl)
	//binds := GenerateBinds(volumes)
	// 配置
	// 处理端口
	for _, port := range ports {
		var bindings []nat.PortBinding
		newPort, _ := nat.NewPort("tcp", strconv.Itoa(port.Number))
		exposePort[newPort] = struct{}{}
		// 判断是internal 还是intranet
		if port.Type == "internal" {
			hostIp = "127.0.0.1"
		} else if port.Type == "intranet" {
			hostIp = "0.0.0.0"
		}
		binding := nat.PortBinding{
			HostIP:   hostIp,
			HostPort: strconv.Itoa(port.Number), // 随机可用端口
		}
		bindings = append(bindings, binding)
		portMap[newPort] = bindings
		portsRsp[port.Number] = port.Number
	}
	// 环境变量
	if len(envs) > 0 {
		for _, env := range envs {
			environments = append(environments, fmt.Sprintf("%s=%s", env.Key, env.Value))
		}
	}
	// host配置
	containerConfig := &container.Config{
		Hostname:     containerName,
		ExposedPorts: exposePort,
		Env:          environments,
		Image:        imageUrlTrim,
		Volumes:      nil,
		//Cmd: []string{"echo", "hello world"},
	}

	var mounts []mount.Mount

	dataDirMount := mount.Mount{
		Type:   mount.TypeBind,
		Source: dataDir,
		Target: "/aospace",
		BindOptions: &mount.BindOptions{
			Propagation:      "",
			NonRecursive:     false,
			CreateMountpoint: false,
		},
	}
	var dockerSockPathOnHost string
	if !strings.Contains(utils.OSKernel(), "microsoft-standard-WSL") {
		dockerSockPathOnHost = "//var/run/docker.sock"
	} else {
		dockerSockPathOnHost = "/var/run/docker.sock"
	}

	sockMount := mount.Mount{
		Type:     mount.TypeBind,
		Source:   dockerSockPathOnHost,
		Target:   "/var/run/docker.sock",
		ReadOnly: true,
	}
	mounts = append(mounts, dataDirMount)
	mounts = append(mounts, sockMount)

	hostConfig := &container.HostConfig{
		//Binds:        binds,
		NetworkMode:  container.NetworkMode(config.Config.Container.Network),
		PortBindings: portMap,
		RestartPolicy: container.RestartPolicy{
			Name:              "always",
			MaximumRetryCount: 0,
		},
		Mounts: mounts,
	}
	// 创建网络
	var networkId string
	networkId, err = ca.DockerAPI.EnsureNetworkExist(config.Config.Container.Network)
	log4go.I("network %s ,networkId %s", config.Config.Container.Network, networkId)
	if networkId == "" {
		networkId, err = ca.DockerAPI.CreateNetwork(containerName)
		if err != nil {
			log4go.E("create network failed")
			return err
		}
	}

	// pull image
	err = ca.DockerAPI.PullImage(imageUrlTrim)
	if err != nil {
		log4go.E("pull image error")
		return err
	}

	exist, err := ca.DockerAPI.EnsureImageExist(imageUrlTrim)
	if err != nil {
		log4go.E("can not find image %s", imageUrlTrim)
		return err
	}
	if exist {
		// 创建前检查容器是否已存在
		cid, err := ca.DockerAPI.FindContainer(containerName)
		if cid == "" {
			log4go.I("container aospace-all-in-one  does not exist")
			containerId, err = ca.DockerAPI.Create(containerConfig, hostConfig, networkConfig, containerName)
			if err != nil {
				log4go.E("create container failed")
				return err
			}
			log4go.D("container aospace-all-in-one is created")
			// 连接网络
			//err = ca.DockerAPI.ConnectNetwork(networkId, containerId)
			//if err != nil {
			//	log4go.I("connect network failed")
			//	return err
			//}
		} else {
			containerStatus, _ := ca.DockerAPI.Inspect(cid)
			if containerStatus.State.Status != "created" {
				return errors.New(fmt.Sprintf("container %s is %s,containerId:%s", containerName, containerStatus.State.Status, containerStatus.ID))
			}
		}
		// 启动容器
		err = ca.DockerAPI.Start(containerId, types.ContainerStartOptions{})
		// 返回
		if err != nil {
			log4go.I("start container failed")
			return err
		}
		// 检查容器状态
		for i := 0; i < 5; i++ {
			containerInfo, _ := ca.DockerAPI.Inspect(containerId)
			if containerInfo.State.Status == "running" {
				log4go.I(fmt.Sprintf("%s status is running", containerName))
				log4go.I("%s upgrade successfully", containerName)
				os.Exit(0)
			} else {
				log4go.I(fmt.Sprintf("%s status is %s,sleep 1 sec", containerName, containerInfo.State.Status))
				time.Sleep(time.Second)
			}
		}
	}
	return fmt.Errorf("not found image %s", imageUrlTrim)

}

func GetAgentRpmPath(versionId string) (string, error) {
	task, err := db.ReadTask(versionId)
	if err != nil {
		return "", fmt.Errorf("GetAgentRpmPath: %s", err)
	}
	if strings.Contains(task.RpmPkg.PkgPath, agentRpmName) && task.RpmPkg.Downloaded {
		return task.RpmPkg.PkgPath, nil
	} else {
		return task.RpmPkg.PkgPath, fmt.Errorf("GetAgentRpmPath: %s", "can't find agent's rpm pkg ")
	}
}

func InstallAgent(versionId string) error {

	rpmPath, err := GetAgentRpmPath(versionId)
	if err != nil {
		log4go.E("Failed to get Agent rpmPath: %s", err)
		fmt.Printf("Failed to get Agent rpmPath: %s", err)
		db.MarkTaskInstallErr(versionId)
		return err
	}
	log4go.I("install system-agent after 5 seconds")
	fmt.Printf("install system-agent after 5 seconds")
	time.Sleep(5 * time.Second)
	err = dnfInstall(rpmPath)
	if err != nil {
		log4go.E("Failed to Install Agent %s", err)
		fmt.Printf("Failed to install Agent: %s", err)
		db.MarkTaskInstallErr(versionId)
		return err
	}
	return nil
}

func GetCurAgentVer() string {
	stdout, _, err := utils.ExeShellCmd("system-agent -version|grep system-agent|awk -F \"-\" '{print $6\"-\"$7}'")
	//log4go.I("stdout:%s", stdout)
	if err != nil {
		return ""
	}
	return stdout
}

func dnfInstall(rpmPath string) error {
	log4go.D("Start to install agent with dnf")
	//_, _, err := utils.ExeCmd("dnf", "update", rpmPath, "-y")
	_, _, err := utils.ExeCmd("rpm", "-Uvh", rpmPath)
	//log4go.I("stdout:%s", stdout)
	if err != nil {
		return fmt.Errorf("dnf install: %v", err)
	}
	log4go.I("Success to install agent with dnf")
	return nil
}

func moveComposeFile(fromPath string, toPath string) error {
	// back raw file
	err := os.Rename(toPath, GetBakPathName(toPath))
	if err != nil {
		return err
	}
	_, err = utils.CopyFile(toPath, fromPath)
	if err != nil {
		return err
	}
	return nil
}

func GetBakPathName(oldPath string) string {
	return oldPath + ".bak"
}

func Install(versionId string) error {
	if versionId != "" {
		oldVer, _ := version.GetInstalledAgentVersion()
		log4go.I("old version eulixspace-agent-%s .", oldVer)
		fmt.Printf("Updating eulixspace-gent-%s...\n", versionId)
		log4go.I("start to install and restart eulixspace-agent-%s .", versionId)

		err := InstallAgent(versionId)
		if err != nil {
			log4go.E("failed to install agent with dnf")
			return err
		}
		curVer := GetCurAgentVer()
		if strings.Trim(curVer, "\n") == versionId {
			log4go.I("current system-agent version %s", curVer)
		}
		task, err := db.ReadTask(versionId)
		if err != nil {
			fmt.Printf("read task info error: %v", err)
			log4go.E("read task info error: %v", err)
			return err
		}
		if task.InstallStatus == models.Done {
			fmt.Printf("upgrade system-agent successfully, current version:%s", versionId)
			log4go.I("upgrade system-agent successfully, current version:%s", versionId)
		}

		//if task.NeedReboot {
		//	log4go.I("upgrade: the aospace server will reboot after 90 seconds")
		//	time.Sleep(90 * time.Second)
		//	utils.ExeCmd("reboot")
		//}
	}
	return nil

}
