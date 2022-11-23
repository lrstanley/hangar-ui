.DEFAULT_GOAL := build-all

export PROJECT := "hangar-ui"
export PACKAGE := "github.com/lrstanley/hangar-ui/cmd/client"
export DOCKER_BUILDKIT := 1

build-all: clean fetch build
	@echo

clean:
	/bin/rm -rfv ${PROJECT}

docker-build:
	docker build \
		--tag ${PROJECT} \
		--force-rm .

license:
	curl -sL https://liam.sh/-/gh/g/license-header.sh | bash -s

# cli
prepare:
	go generate -x ./...

fetch:
	go mod download
	go mod tidy

upgrade-deps:
	go get -u ./...
	go mod tidy

upgrade-deps-patch:
	go get -u=patch ./...
	go mod tidy

dlv: prepare
	dlv debug \
		--headless --listen=:2345 \
		--api-version=2 \
		${PACKAGE} -- --debug

debug: prepare
	go run ${PACKAGE} --debug --log.path debug.log --target ci

build: prepare fetch
	CGO_ENABLED=0 \
	go build \
		-ldflags '-d -s -w -extldflags=-static' \
		-tags=netgo,osusergo,static_build \
		-installsuffix netgo \
		-buildvcs=false \
		-trimpath \
		-o ${PROJECT} \
		${PACKAGE}
