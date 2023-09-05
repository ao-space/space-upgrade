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

package models

import (
	"reflect"
	"time"
)

const (
	Ing         = "ing"
	Done        = "done"
	Err         = "err"
	Downloading = "downloading"
	Downloaded  = "downloaded"
	Installing  = "installing"
	Installed   = "installed"
	DownloadErr = "download-err"
	InstallErr  = "install-err"
)

type Task struct {
	VersionId        string          `json:"versionId"`
	Status           string          `json:"status"`        // 整体流程状态："", downloading, downloaded, installing, installed, download-err，install-err
	DownStatus       string          `json:"downStatus"`    // 下载状态："", ing, done, err
	InstallStatus    string          `json:"installStatus"` // 安装状态："", ing, done, err
	StartDownTime    string          `json:"startDownTime"`
	StartInstallTime string          `json:"startInstallTime"`
	DoneDownTime     string          `json:"doneDownTime"`
	DoneInstallTime  string          `json:"doneInstallTime"`
	RpmPkg           VersionDownInfo `json:"rpmPkg"`
	CFile            VersionDownInfo `json:"cFile"` // docker-compose.yml
	ContainerImg     VersionDownInfo `json:"containerImg"`
	NeedReboot       bool            `json:"reboot"`
}

func (t *Task) MarkInstallErr() {
	t.Status = InstallErr
	t.InstallStatus = Err
	t.DoneInstallTime = time.Now().Format(time.RFC3339)
}

type VersionDownInfo struct {
	VersionId  string    `json:"versionId"`
	Downloaded bool      `json:"downloaded"`
	PkgPath    string    `json:"pkgPath"`
	UpdateTime time.Time `json:"updateTime"`
}

type TimeTransformer struct {
}

type StartUpgradeReq struct {
	VersionId string `json:"versionId" form:"versionId"`
	DataDir   string `json:"dataDir" form:"dataDir"`
}

func (t TimeTransformer) Transformer(typ reflect.Type) func(dst, src reflect.Value) error {
	if typ == reflect.TypeOf(time.Time{}) {
		return func(dst, src reflect.Value) error {
			if dst.CanSet() {
				isZero := dst.MethodByName("IsZero")
				result := isZero.Call([]reflect.Value{})
				if result[0].Bool() {
					dst.Set(src)
				}
			}
			return nil
		}
	}
	return nil
}
