.PHONY: all tags clean test build install generate image release

APP_NAME=service-api
BIN_NAME=service-api
APP_REPO_URL="github.com/sandbreaker/goservice"

OK_COLOR=\033[32;01m
NO_COLOR=\033[0m
ERROR_COLOR=\033[31;01m
WARN_COLOR=\033[33;01m

# Build Flags
BUILD_DATE = $(shell date -u)
BUILD_HASH = $(shell git rev-parse --short HEAD)
BUILD_NUMBER ?= $(BUILD_NUMBER:)
VERSION = 0.1.1-$(BUILD_HASH)

# If we don't set the build number it defaults to dev
ifeq ($(BUILD_NUMBER),)
	BUILD_NUMBER := dev
endif

NOW = $(shell date -u '+%Y%m%d%I%M%S')

DOCKER := docker
GO := go
GO_ENV := $(shell $(GO) env GOOS GOARCH)
GOOS ?= $(word 1,$(GO_ENV))
GOARCH ?= $(word 2,$(GO_ENV))
GOFLAGS ?= $(GOFLAGS:)
ROOT_DIR := $(realpath .)

GOSERVICE_ENV ?= "dev"

# GOOS/GOARCH of the build host, used to determine whether 
# we're cross-compiling or not
BUILDER_GOOS_GOARCH="$(GOOS)_$(GOARCH)"

PKGS = $(shell $(GO) list ./...)
# PKGS = $(shell $(GO) list ./... ./cmd/... ./pkg/... | grep -v /vendor/)

TAGS ?= "netgo"
BUILD_ENV =
ENVFLAGS = CGO_ENABLED=1 $(BUILD_ENV)

ifneq ($(GOOS), darwin)
	EXTLDFLAGS = -extldflags "-lm -lstdc++ -static"
else
	EXTLDFLAGS =
endif

GO_LINKER_FLAGS ?= --ldflags \
	'$(EXTLDFLAGS) -s -w -X "$(APP_REPO_URL)/version.BuildNumber=$(BUILD_NUMBER)" \
   -X "$(APP_REPO_URL)/version.BuildDate=$(BUILD_DATE)" \
   -X "$(APP_REPO_URL)/version.BuildHash=$(BUILD_HASH)"'

all: build

tidy:
	@echo "Tidies go modules update/add..."
	@$(GO) mod tidy

release:
	@echo "Generating release manifest file for $(VERSION)"
	echo $(VERSION) > .cfg/version

generate:
	@echo "$(OK_COLOR)==> Generating files via go generate...$(NO_COLOR)"
	@$(GO) generate $(GOFLAGS) $(PKGS)

# build-static: test
# 	@echo "$(OK_COLOR)==> Building binary ($(GOOS)/$(GOARCH))...$(NO_COLOR)"
# 	@echo @$(ENVFLAGS) $(GO) build -a -installsuffix cgo -tags $(TAGS)  $(GOFLAGS) $(GO_LINKER_FLAGS) -o bin/$(GOOS)_$(GOARCH)/$(BIN_NAME) cmd/api/main.go
# 	@$(ENVFLAGS) $(GO) build -a -installsuffix cgo -tags $(TAGS) $(GOFLAGS) $(GO_LINKER_FLAGS) -o bin/$(GOOS)_$(GOARCH)/$(BIN_NAME) cmd/api/main.go

# install-static: build
# 	@echo "$(OK_COLOR)==> Installing packages into GOPATH...$(NO_COLOR)"
# 	@$(GO) install $(GOFLAGS) $(PKGS)

install:
	@$(GO) install ./...

run-dev: install
	GOSERVICE_ENV=$(GOSERVICE_ENV) GOSERVICE_DEBUG=true $(GOPATH)/bin/$(BIN_NAME)

format:
	@$(GO) fmt ./... && $(GO) vet ./...


##########
# HELPERS!
##########

# generate script to login to aws docker repo
CMD_REPOLOGIN := "eval $$\( aws ecr"
ifdef AWS_CLI_PROFILE
CMD_REPOLOGIN += " --profile $(AWS_CLI_PROFILE)"
endif
ifdef AWS_CLI_REGION
CMD_REPOLOGIN += " --region $(AWS_CLI_REGION)"
endif
CMD_REPOLOGIN += " get-login --no-include-email \)"

# login to AWS-ECR
repo-login: ## Auto login to AWS-ECR unsing aws-cli
	@eval $(CMD_REPOLOGIN)

