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

package config

import (
	"fmt"
	"os"
	"time"

	"testing"

	"github.com/dungeonsnd/gocom/file/fileutil"
	"github.com/jinzhu/configor"
)

const (
	appName = "eulixspace-upgrade"
)

var Config = struct {
	Log struct {
		Path          string `default:"/opt/logs/eulixspace-upgrade"`
		Filename      string `default:"upgrade"`
		MaxAge        uint32 `default:"0"`        // 文件最多保留多少小时后被覆盖. MaxAge和RotationTime一起使用. MaxAge不可与RotationSize一起使用.
		RotationTime  int64  `default:"0"`        // 多少小时后生成新文件. 小于等于0表示禁用.
		RotationSize  int64  `default:"10485760"` // 多少Byte生成新文件. 小于等于0表示禁用. RotationSize和RotationCount一起使用. MaxAge不可与RotationSize一起使用.
		RotationCount uint32 `default:"20"`       // 多少个文件之后覆盖最早的文件.
		Level         uint32 `default:"5"`        // InfoLevel=4, DebugLevel=5, logrus.TraceLevel=6
		Dsn           string `default:"https://7856d40232e44eaabd0ee539595a95dd@sentry.eulix.xyz/17"`
	}

	Settings struct {
		AutoDownload bool `default:"true"`
		AutoInstall  bool `default:"true"`
	}
	Docker struct {
		APIVersion  string `default:"1.39"`
		NetworkName string `default:"bp-cicada"`
		Registry    struct {
			UserName string `default:"636963616461"`
			Password string `default:"5a4a436f6f563767455f69563538722d32764d32"`
			Url      string `default:"registry.eulix.xyz"`
		}
	}

	Platform struct {
		APIBase struct {
			Url string `default:"https://services.ao.space/platform"`
		}

		LatestVersion struct {
			Path string `default:"v1/api/package/box"`
		}
	}

	UpgradeConfig struct {
		SettingsFile string `default:"/etc/bp/upgrade/settings.json"`
		// AutoDownload bool   `default:"true" json:"autoDownload"`
		// AutoInstall  bool   `default:"true" json:"autoInstall"`
	}

	RunTime struct {
		BasePath          string `default:"/var/system-agent/"`
		DBDir             string `default:".db"`
		PkgDir            string `default:"pkg"`
		UpgradeCollection string `default:"upgrade"`
		TaskResource      string `default:"task"`
		IpAddr            string `default:"172.17.0.1:5681"`
	}
	Container struct {
		Name            string `default:"aospace-all-in-one"`
		ImageOpenSource string `default:"hub.eulix.xyz/ao-space/space-agent"`
		Restart         bool   `default:"true"`
		Network         string `default:"ao-space"`
		WebPort         int    `default:"5678"`
		ApiPort         int    `default:"5680"`
		Env             map[string]interface{}
		Volume          []string
	}
}{}

var ConfFile = "/etc/bp/" + appName + ".yml"

func init() {
	testing.Init()

	err := configor.New(&configor.Config{AutoReload: true,
		AutoReloadInterval: time.Second * 15,
		AutoReloadCallback: func(config interface{}) {
			fmt.Printf("config file changed:\n%+v\n", config)
		}}).Load(&Config, ConfFile)
	if err != nil {
		fmt.Printf("Failed to load config file: %v", err)
	}

	createLogFileDir(Config.Log.Path)
	//writeDefaultConfigFile(ConfFile)

	fmt.Printf("config: %+v\n\n", Config)

	if RunInContainer() {
		Config.RunTime.BasePath = "/aospace/var/system-agent"
	}
	if _, err = os.Stat(Config.RunTime.BasePath); err != nil {
		os.Mkdir(Config.RunTime.BasePath, os.ModePerm)
	}
}

func createLogFileDir(dir string) {
	if !fileutil.IsFileExist(dir) {
		err := fileutil.CreateDirRecursive(dir)
		if err != nil {
			fmt.Printf("Failed to create log directory %s: %v \n", dir, err)
		}
	}
}

func RunInContainer() bool {
	if _, err := os.Stat("/.dockerenv"); err == nil {
		return true
	}
	return false
}

//func writeDefaultConfigFile(f string) {
//
//	out, err := yaml.Marshal(Config)
//	if err != nil {
//		fmt.Printf("Failed  yaml.Marshal: %+v\n", err)
//		return
//	}
//	err = fileutil.WriteToFile(f, out, true)
//	if err != nil {
//		log4go.E("Failed to write Default ConfigFile: %v", err)
//	}
//}
