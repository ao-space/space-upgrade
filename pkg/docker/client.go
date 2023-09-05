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
	"github.com/docker/docker/client"
	"log"
	"net"
)

const (
	ClientVersion = "1.39"
	SocketAddr    = "/var/run/docker.sock"
)

func DockerClient() *client.Client {
	dockerClient, err := client.NewClientWithOpts(
		//client.WithHost("tcp://192.168.124.84:2375"),
		client.WithDialContext(func(_ context.Context, _, _ string) (net.Conn, error) {
			return net.Dial("unix", SocketAddr)
		}),
		client.WithVersion(ClientVersion))
	if err != nil {
		log.Fatal(err)
	}
	ping, err := dockerClient.Ping(context.Background())
	//}))
	if err != nil {
		log.Fatal(err)
	}
	log.Println(ping.APIVersion)
	return dockerClient
}
