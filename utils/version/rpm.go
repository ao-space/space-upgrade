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

package version

import (
	"eulixspace-upgrade/utils"
	"fmt"
	"strings"
)

// GetInstalledAgentVersion from rpm
func GetInstalledAgentVersion() (string, error) {
	queryCmd := fmt.Sprintf(`dnf list installed %s | grep %s.aarch64  | awk '{print $2}'`, "eulixspace-agent", "eulixspace-agent")
	stdout, stderr, err := utils.ExeCmd("bash", "-c", queryCmd)
	if err != nil || stderr != "" {
		return stdout, fmt.Errorf("get agent version though rpm %s: %s", err, stderr)
	}
	versionId := strings.TrimSpace(stdout)
	return versionId, nil
}

//func CleanAllAndMakeCache() error {
//	cmd := `dnf clean all && dnf makecache`
//	_, stderr, err := utils.ExeShellCmd(cmd)
//	if err != nil || stderr != "" {
//		return fmt.Errorf("dnf clean all and make cache error: %v: %v", err, stderr)
//	}
//	return nil
//}
