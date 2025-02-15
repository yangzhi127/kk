.PHONY: all build docker

BINARY_NAME=awsKeyHunter
DOCKER_IMAGE_NAME=aws-key-hunter:latest

all: build docker

build:
	go build -o $(BINARY_NAME) cmd/awsKeyhunter.go

docker:
	sudo docker build -t $(DOCKER_IMAGE_NAME) .