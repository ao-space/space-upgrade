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

package cmd

import (
	"eulixspace-upgrade/models"
	"eulixspace-upgrade/utils"
	"eulixspace-upgrade/utils/db"
	"eulixspace-upgrade/utils/version"
	"eulixspace-upgrade/views"
	"fmt"
	"github.com/dungeonsnd/gocom/log4go"
	"github.com/spf13/cobra"
	"strings"
	"time"
)

//var Wg sync.WaitGroup

var RootCmd = &cobra.Command{
	Use:   "eulixspace-upgrade",
	Short: "A command line tool to upgrade system and Docker containers",
}

var Update = &cobra.Command{
	Use:   "install",
	Short: "install newer system-agent",
	Run: func(cmd *cobra.Command, args []string) {
		versionId, _ := cmd.Flags().GetString("version")
		if versionId != "" {
			oldVer, _ := version.GetInstalledAgentVersion()
			log4go.I("old version eulixspace-agent-%s .", oldVer)
			fmt.Printf("Updating eulixspace-gent-%s...\n", versionId)
			log4go.I("start to install and restart eulixspace-agent-%s .", versionId)
			//Wg.Add(1)
			//done := make(chan bool, 1)
			//quit := make(chan os.Signal, 1)
			//signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
			go func() {
				//defer Wg.Done()
				err := views.InstallAgent(versionId)
				if err != nil {
					fmt.Printf("install agent error:%v", err)
				}
				//done <- true
			}()
			//select {
			//case <-done:
			//	fmt.Println("dnf update system-agent command has been executed successfully.")
			//case <-quit:
			//	fmt.Println("Interrupt signal received. Aborting...")
			//	os.Exit(1)
			//}
			for i := 1; i <= 120; i++ {
				curVer := views.GetCurAgentVer()
				log4go.I("loop %d: current system-agent version %s", i, curVer)
				if strings.Trim(curVer, "\n") == versionId {
					break
				}
				time.Sleep(time.Second)
			}
			_, err := db.MarkTaskInstalled(versionId)
			if err != nil {
				fmt.Printf("mark task installed error:%v", err)
			}

		}
		task, err := db.ReadTask(versionId)
		if err != nil {
			fmt.Printf("read task info error: %v", err)
			log4go.E("read task info error: %v", err)
		}
		if task.InstallStatus == models.Done {
			fmt.Printf("upgrade system-agent successfully, current version:%s", versionId)
			log4go.I("upgrade system-agent successfully, current version:%s", versionId)
		}

		if task.NeedReboot {
			log4go.I("upgrade: the aospace server will reboot after 10 seconds")
			time.Sleep(10 * time.Second)
			utils.ExeCmd("reboot")
		}
		//fmt.Println("Docker containers upgrade completed successfully!")
	},
}
