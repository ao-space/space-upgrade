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
	"eulixspace-upgrade/models"
	"github.com/dungeonsnd/gocom/log4go"
	"github.com/gin-gonic/gin"
	"net/http"
)

// StartUpgrade godoc
// @Summary POST:/upgrade/start-upgrade 开始安装最新版本
// @Tags upgrade
// @Accept   json
// @Produce   json
// @Param upgrade body models.StartUpgradeRes true "安装"
// @Success 200 {object} models.Task
// @Failure 400 string models.BaseRsp  "版本尚未下载，请先下载"
// @Failure 409 string models.BaseRsp  "已经有一个任务正在启动"
// @Failure 500 string models.BaseRsp  "失败，稍后再试"
// @Router /upgrade/v1/api/start [POST]
func StartUpgrade(c *gin.Context) {
	log4go.I("/agent/v1/api/upgrade/start")

	var req models.StartUpgradeReq
	err := c.ShouldBindJSON(&req)
	if err != nil {
		log4go.E("Failed to bind args %s", err)
		c.JSON(http.StatusBadRequest, models.BaseRsp{Code: "UP-400", Message: err.Error()})
		return
	}
	containerAgent := NewContainerAgent()
	//go InstallAgent(task.VersionId)
	go containerAgent.Upgrade(req.VersionId, req.DataDir)
	c.JSON(http.StatusOK, models.BaseRsp{
		Code:    "UP-200",
		Message: "ok",
	})
}
