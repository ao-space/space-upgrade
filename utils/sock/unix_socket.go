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

package sock

import (
	"context"
	"eulixspace-upgrade/views"
	"github.com/dungeonsnd/gocom/log4go"
	"net"
)

const SockPath = "/var/system-agent/upgrade.sock"

func HandleUnixSocket(ctx context.Context, conn *net.UnixConn) {
	agent := &views.NativeAgent{}
	select {
	case <-ctx.Done():
		log4go.I("Received signal. Shutting down gracefully...")
		return
	default:
		//conn, err := l.Accept()
		//if err != nil {
		//	log4go.E("Failed to accept incoming connection: %v", err)
		//}

		// 处理来自客户端的请求
		buf := make([]byte, 1024)
		n, err := conn.Read(buf)
		if err != nil {
			log4go.E("Failed to read data from UNIX socket: %v", err)
		}
		log4go.I("Received %d bytes of data from client: %s", n, string(buf[:n]))
		conn.Write([]byte("OK"))
		agent.Upgrade(string(buf[:n]))
	}

}
