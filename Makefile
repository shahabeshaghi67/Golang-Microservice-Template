NAME				:= golang-api-service
BINARY				:= $(NAME)
VERSION				:= 
DOCKER_IMAGE		:= $(BINARY)
CONTAINER_REPO		:= 
API_GENERATOR_DIR	:= ./cmd/api-doc-generator

GIT_HEAD         	:= $(shell git rev-parse HEAD)
PACKAGES         	:= $(shell find . -name *.go | grep -v -e vendor | xargs -n1 dirname | sort -u)
MAIN_DIR         	:= ./cmd/$(NAME)
TEST_FLAGS       	:= -v -race -count=1 -mod=readonly -cover -coverprofile coverprofile.txt
LINK_FLAGS       	:= -X main.Version=$(VERSION) -X main.GitHead=$(GIT_HEAD)
BUILD_FLAGS      	:= -mod=readonly -v
DEP_BASE_VERSION 	?= latest
NO_INTEGRATION   	?= 0
ifeq ($(NO_INTEGRATION),0)
TEST_FLAGS += -tags=integration
endif

default: build

generate:
	go generate $(PACKAGES)

api-doc: generate
	go mod tidy
	CGO_ENABLED=0 go run $(API_GENERATOR_DIR)

test: generate
	go mod tidy
	golangci-lint --version
	golangci-lint run $(PACKAGES)
	# go test $(TEST_FLAGS) $(PACKAGES)

build: test 
	CGO_ENABLED=0 go build $(BUILD_FLAGS) -ldflags="$(LINK_FLAGS)" -o build/$(NAME) $(MAIN_DIR)
	@echo build complete

build.linux: build-service.linux 
	@echo builds complete

build-service.linux: test
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build $(BUILD_FLAGS) -ldflags="$(LINK_FLAGS)" -o build/linux/$(NAME) $(MAIN_DIR)
	@echo build-service.linux complete

docker: docker-service 
	docker push $(CONTAINER_REPO)/$(DOCKER_IMAGE):$(VERSION)
	@echo docker builds complete and pushed

docker-service: build-service.linux
	docker build -t $(CONTAINER_REPO)/$(DOCKER_IMAGE):$(VERSION) -t $(CONTAINER_REPO)/$(DOCKER_IMAGE):latest \
	--build-arg DEP_BASE_VERSION=${DEP_BASE_VERSION} CONTAINER_REPO=${CONTAINER_REPO} .
	@echo docker build of image $(CONTAINER_REPO)/$(DOCKER_IMAGE):$(VERSION) complete

clean:
	rm -rvf pkg/mocks build coverprofile.txt

.PHONY: clean
