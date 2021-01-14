GOPATH:=$(shell go env GOPATH)

# ----- Installing -----

.PHONY: install
install:
	go mod download

.PHONY: lint
lint:
	@ [ -e ./bin/golangci-lint ] || wget -O - -q https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s latest
	./bin/golangci-lint run

# ----- Testing -----

BUILDENV := CGO_ENABLED=0
TESTFLAGS := -short -cover
SERVICE ?= $(base_dir)

.PHONY: test
test:
	$(BUILDENV) go test $(TESTFLAGS) ./...


.PHONY: build
build:
	CGO_ENABLED=0 go build $(SERVICE)


.PHONY: dev
dev:
	go run .

DOCKER_REGISTRY=freemiumvpn
DOCKER_CONTAINER_NAME=fpn-auth
DOCKER_REPOSITORY=$(DOCKER_REGISTRY)/$(DOCKER_CONTAINER_NAME)
SHA8=$(shell echo $(GITHUB_SHA) | cut -c1-8)

docker-image:
	docker build --rm \
		--tag $(DOCKER_REPOSITORY):local . \
		--build-arg SERVICE=$(SERVICE)

ci-docker-auth:
	@echo "${DOCKER_PASSWORD}" | docker login --username "${DOCKER_USERNAME}" --password-stdin

ci-docker-build:
	@docker build --rm \
		--tag $(DOCKER_REPOSITORY):$(SHA8) \
		--tag $(DOCKER_REPOSITORY):latest .

ci-docker-build-push: ci-docker-build
	@docker push $(DOCKER_REPOSITORY):$(SHA8)
	@docker push $(DOCKER_REPOSITORY):latest
