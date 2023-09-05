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

package utils

import (
	"bytes"
	"fmt"
	"github.com/dungeonsnd/gocom/encrypt/hash/sha256"
	"github.com/dungeonsnd/gocom/file/fileutil"
	"github.com/dungeonsnd/gocom/log4go"
	"io"
	"os"
	"os/exec"
	"regexp"
	"strings"
)

func ExeShellCmd(queryCmd string) (string, string, error) {
	return ExeCmd("bash", "-c", queryCmd)
}

func ExeCmd(name string, arg ...string) (string, string, error) {
	log4go.D("exec: %s %s", name, arg)
	cmd := exec.Command(name, arg...)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
	outStr, errStr := string(stdout.Bytes()), string(stderr.Bytes())
	return outStr, errStr, err
}

func CopyFile(dstName, srcName string) (written int64, err error) {
	src, err := os.Open(srcName)
	if err != nil {
		return
	}
	defer src.Close()

	dst, err := os.OpenFile(dstName, os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return
	}
	defer dst.Close()

	return io.Copy(dst, src)
}

//func IpAddrExist(ip string) bool {
//	if strings.Contains(ip, ":") {
//		ip = strings.Split(ip, ":")[0]
//	}
//	ip = strings.TrimSpace(ip)
//
//	_, stderr, err := ExeShellCmd(fmt.Sprintf(`ip add | grep "%s"`, ip))
//	if err != nil || stderr != "" {
//		return false
//	}
//	return true
//}

func GetRPIBtId() (string, error) {
	b, err := fileutil.ReadFromFile("/proc/cpuinfo")
	if err != nil {
		log4go.W("@@@@ GetBtId, failed ReadFromFile, err:%v", err)
		return "", fmt.Errorf("failed ReadFromFile, err:%v", err)
	}
	s := string(b)

	r, err := regexp.Compile(".*Serial.*[a-zA-Z0-9]*")
	if err != nil {
		return "", fmt.Errorf("failed regexp.Compile r, err:%v", err)
	}
	ser := r.FindString(s)
	ser = strings.Replace(ser, "Serial", "", -1)
	ser = strings.Replace(ser, ":", "", -1)
	ser = strings.Replace(ser, "\t", "", -1)
	ser = strings.Replace(ser, " ", "", -1)
	ser = strings.Replace(ser, "\r", "", -1)
	ser = strings.Replace(ser, "\n", "", -1)
	serExtend := fmt.Sprintf("eulixspace-btid-%v", ser)
	btid := sha256.HashHex([]byte(serExtend), 1)[:16]
	return btid, nil
}

func OSKernel() string {
	content, err := os.ReadFile("/proc/version")
	if err != nil {
		return ""
	}
	kernelVer := strings.Replace(string(content), "\n", "", -1)
	return kernelVer
}

func RunInContainer() bool {
	if _, err := os.Stat("/.dockerenv"); err != nil {
		return false
	}
	return true
}
