# Copyright 2016 The Kubernetes Authors.
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

# This version-strategy uses git tags to set the version string
VERSION ?= $(shell git describe --tags --always --dirty)
#
# This version-strategy uses a manual value to set the version string
#VERSION ?= 0.0.1

help: # @HELP 打印帮助信息
help:
	@echo "VARIABLES:"
	@echo
	@echo "TARGETS:"
	@grep -E '^.*: *# *@HELP' $(MAKEFILE_LIST)    \
	    | awk '                                   \
	        BEGIN {FS = ": *# *@HELP"};           \
	        { printf "  %-30s %s\n", $$1, $$2 };  \
		'

version: # @HELP 版本信息
version:
	@echo $(VERSION)

mock: # @HELP 为 Repository 接口生成 mock 代码
mock:
	mockgen -source=dal/user/repository.go -destination=dal/user/mock.go -package=user
	mockgen -source=dal/bill/repository.go -destination=dal/bill/mock.go -package=bill
	mockgen -source=dal/telegram/repository.go -destination=dal/telegram/mock.go -package=telegram
	mockgen -source=telebot/user_state.go -destination=telebot/user_state_mock.go -package=telebot
	mockgen -destination=mock/telebotmock/mock.go -package=telebotmock gopkg.in/telebot.v3 Context
	@echo Done.

test: # @HELP 运行单元测试
test:
	@GOARCH=amd64 go test -gcflags=all=-l -cover \
		./service/...\
		./conf/...\
		./telebot/...

.PHONY: help version mock test