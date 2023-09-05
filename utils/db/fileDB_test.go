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
	"eulixspace-upgrade/models"
	"github.com/stretchr/testify/assert"
	"testing"
)

func InItEmptyTask() (*models.Task, error) {
	emptyTask := models.Task{}
	return UpdateOrCreateTask(&emptyTask)
}

func TestUpdateOrCreateTaskInItTask(t *testing.T) {
	emptyTask := models.Task{}
	newTask, err := InItEmptyTask()
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, newTask.VersionId, emptyTask.VersionId)
	assert.Equal(t, newTask.Status, emptyTask.Status)
}
