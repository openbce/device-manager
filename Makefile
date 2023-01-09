IMG_NAME=openbce/device-manager
IMG_VER=v0.1


device-manager: init
	go build -o _output/device-manager cmd/device-manager/main.go

init:
	mkdir -p _output

docker-build:
	docker build -t ${IMG_NAME}:${IMG_VER} .
