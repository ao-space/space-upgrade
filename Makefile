# Copyright (c) 2022 Institute of Software, Chinese Academy of Sciences (ISCAS)
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

NAME=eulixspace-upgrade
#OPTIONS=-trimpath -mod=vendor
LDFLAGS=-ldflags "-s -w"
SOURCES=$(shell ls **/*.go)
# CHECKS:=check

.PHONY: all
all: exe

.PHONY: exe
exe: $(SOURCES) Makefile
	echo "building..."
	go env -w GO111MODULE=on
	go get github.com/pkg/errors
	go get github.com/dungeonsnd/gocom/...
	GOOS=linux GOARCH=${ARCH} go build $(OPTIONS) $(LDFLAGS) -o build/$(NAME)

.PHONY: clean
clean:
	go clean -i
	rm -rf build/*
