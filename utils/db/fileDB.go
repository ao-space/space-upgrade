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

package db

import (
	"eulixspace-upgrade/config"
	"eulixspace-upgrade/models"
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/dungeonsnd/gocom/log4go"

	"github.com/imdario/mergo"

	"path"

	scribble "github.com/nanobox-io/golang-scribble"
)

// 这里可以存在客户端并发请求,导致并发读写文件. 故尝试在此文件读写处增加锁.
var lock sync.RWMutex

var conf = config.Config.RunTime
var Dir = path.Join(conf.BasePath, conf.DBDir)

func NewDBClient() (*scribble.Driver, error) {
	return scribble.New(Dir, nil)
}

func initDB(filePath string) error {
	defer lock.Unlock()
	lock.Lock()
	f, err := os.OpenFile(filePath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	_, err = f.Write([]byte("{}"))
	if err != nil {
		return err
	}
	f.Sync()
	defer f.Close()
	return nil
}

// UpdateOrCreateTask to write db for create a new task or update old task, according to versionId.
func UpdateOrCreateTask(newT *models.Task) (*models.Task, error) {
	defer lock.Unlock()
	lock.Lock()
	task := new(models.Task)
	db, err := NewDBClient()
	if err != nil {
		return task, err
	}
	err = db.Read(conf.UpgradeCollection, conf.TaskResource, &task)
	if err != nil {
		return task, err
	}

	if newT.VersionId != task.VersionId {
		// 一个新的 Task， 需要覆盖空值
		err = db.Write(conf.UpgradeCollection, conf.TaskResource, newT)
		if err != nil {
			return task, fmt.Errorf("update task => %w", err)
		}
		return newT, nil

	} else {
		// 相同的 task，空值不覆盖
		err = mergo.Merge(task, newT, mergo.WithTransformers(models.TimeTransformer{}), mergo.WithOverride)
		if err != nil {
			return task, fmt.Errorf("update task => %w", err)
		}
		err = db.Write(conf.UpgradeCollection, conf.TaskResource, task)
		if err != nil {
			return task, fmt.Errorf("update task => %w", err)
		}
		return task, nil
	}
}

// ReadTask is to read a task from the database, when versionId is set, It will match the version.
func ReadTask(versionId string) (*models.Task, error) {
	defer lock.RUnlock()
	lock.RLock()
	task := new(models.Task)
	db, err := NewDBClient()
	if err != nil {
		return task, err
	}
	err = db.Read(conf.UpgradeCollection, conf.TaskResource, &task)
	if err != nil {
		return task, err
	}
	if versionId != "" && versionId != task.VersionId {
		return task, fmt.Errorf("no record of the specified version exists")
	}
	return task, nil
}

func MarkTaskInstalled(versionId string) (*models.Task, error) {
	log4go.D("Marking task %s status %s", versionId, models.Installed)
	doc, err := UpdateOrCreateTask(&models.Task{
		VersionId:     versionId,
		Status:        models.Installed,
		InstallStatus: models.Done,
		DoneDownTime:  time.Now().Format(time.RFC3339)})
	return doc, err
}

func MarkTaskInstallErr(versionId string) *models.Task {
	log4go.D("Marking task %s status %s", versionId, models.InstallErr)
	doc, err := UpdateOrCreateTask(&models.Task{
		VersionId:       versionId,
		Status:          models.InstallErr,
		InstallStatus:   models.Err,
		DoneInstallTime: time.Now().Format(time.RFC3339)})
	if err != nil {
		log4go.E("Failed to mark task up err %s", err)
	}
	return doc
}
