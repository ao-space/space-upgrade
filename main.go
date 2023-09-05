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

package main

import (
	"context"
	"eulixspace-upgrade/config"
	"eulixspace-upgrade/info"
	"eulixspace-upgrade/utils"
	"eulixspace-upgrade/utils/sock"
	"eulixspace-upgrade/views"
	"github.com/dungeonsnd/gocom/log4go"
	"github.com/evalphobia/logrus_sentry"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
)

func initRouter() *gin.Engine {
	router := gin.Default()
	router.GET("/ping", func(c *gin.Context) {
		c.String(200, "pong")
	})
	upgradeRoute := router.Group("/upgrade/v1/api")
	{
		upgradeRoute.POST("/start", views.StartUpgrade)
	}
	return router
}

func main() {
	log4go.InitLog(config.Config.Log.Path, config.Config.Log.Filename,
		config.Config.Log.MaxAge, config.Config.Log.RotationTime,
		config.Config.Log.RotationSize, config.Config.Log.RotationCount)
	log4go.SetLogLevel(5) // InfoLevel=4, DebugLevel=5, logrus.TraceLevel=6

	log4go.I("================[Started] [eulixspace-upgrade]================")
	btid, err := utils.GetRPIBtId()
	if err != nil {
		log4go.E("Failed to get RPI BtId: %v", err)
	}
	hook, err := logrus_sentry.NewAsyncWithTagsSentryHook(
		config.Config.Log.Dsn,
		map[string]string{"version": info.Version, "btid": btid},
		[]logrus.Level{
			logrus.PanicLevel,
			logrus.FatalLevel,
			logrus.ErrorLevel,
			logrus.WarnLevel,
			logrus.InfoLevel,
		})
	if err != nil {
		log.Println(err)
	}
	log4go.BindHooks(hook)

	// 创建Unix套接字
	socketPath := sock.SockPath
	if err := os.RemoveAll(socketPath); err != nil {
		log.Fatalf("Failed to remove socket file %q: %v", socketPath, err)
	}

	l, err := net.ListenUnix("unix", &net.UnixAddr{Name: socketPath, Net: "unix"})
	if err != nil {
		log4go.E("Failed to listen on UNIX socket %q: %v", socketPath, err)
	}
	defer l.Close()

	log4go.I("Listening on UNIX socket %q", socketPath)

	newRouter := initRouter()
	if utils.RunInContainer() {
		err = newRouter.Run(":5681")
		if err != nil {
			log4go.F("run http web server err:%v", err)
		}
	}

	// 监听信号并创建context
	ctx, cancel := context.WithCancel(context.Background())
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		sig := <-signals
		log4go.I("Received signal %v. Shutting down gracefully...", sig)

		// 取消context，通知goroutine退出
		cancel()
		os.Exit(0)
	}()

	// 实时监听Unix套接字
	for {
		conn, err := l.AcceptUnix()
		if err != nil {
			log.Fatalf("Failed to accept incoming connection: %v", err)
		}

		// 处理连接
		go sock.HandleUnixSocket(ctx, conn)
	}

}
