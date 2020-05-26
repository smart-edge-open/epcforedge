# SPDX-License-Identifier: Apache-2.0
# Copyright (c) 2019-2020 Intel Corporation

# Source user configured environment file.
include .env

export GO111MODULE = on

define bring_ui_up
	docker-compose up -d cups-ui
	docker-compose up -d cnca-ui
endef

define bring_ui_down
	docker-compose stop cups-ui
	docker-compose stop cnca-ui
endef

.PHONY: help clean build lint ui-up ui-down

help:
	@echo "Please use \`make <target>\` where <target> is one of"
	@echo "Building:"
	@echo "  clean             to clean up build artifacts and docker volumes"
	@echo "  build             to build the project to the ./dist/ folder"
	@echo ""
	@echo "Services:"
	@echo "  ui-up             to start all of the production UI Containers"
	@echo "  ui-down           to stop all of the production UI containers"
	@echo ""
	@echo "  cups-ui-up        to start the production UI Container"
	@echo "  cups-ui-down      to stop the production UI container"
	@echo "  cups-ui-dev-up    to start local developer instance of the UI"
	@echo "  cups-ui-test      run the UI project tests"
	@echo ""
	@echo "  cnca-ui-up        to start the production UI Container"
	@echo "  cnca-ui-down      to stop the production UI container"
	@echo "  cnca-ui-dev-up    to start local developer instance of the UI"
	@echo "  cnca-ui-test      run the UI project tests"

clean:
	@docker-compose stop
	@docker-compose rm
	rm -rf dist certificates artifacts

build:
	docker-compose build cups-ui cnca-ui

lint:
	golangci-lint run

ui-up:
	$(call bring_ui_up)

ui-down:
	$(call bring_ui_down)

cups-ui-up:
	docker-compose up -d cups-ui

cups-ui-down:
	docker-compose stop cups-ui

cups-ui-dev-up:
	cd cups && yarn install && yarn start

cups-ui-test:
	cd cups && yarn install && yarn build && yarn test

cnca-ui-up:
	docker-compose up -d cnca-ui

cnca-ui-down:
	docker-compose stop cnca-ui

cnca-ui-dev-up:
	cd cnca && yarn install && yarn start

cnca-ui-test:
	cd cnca && yarn install && yarn build && yarn test
