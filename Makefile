NAME=kube-debug
VERSION=v0.1.0
NET_MODE=host
IMAGE_TAG=cloudnativer/kube-debug

all:kube-debug kube-debug-image

kube-debug:
	@echo Start building kube-debug.
	go build -o $(NAME) kube-debug.go
	@echo Finished building.

kube-debug-image:
	@echo Start making container image.
	docker build --network=$(NET_MODE) -t $(IMAGE_TAG):$(VERSION) .
	docker save $(IMAGE_TAG):$(VERSION) -o $(NAME)-container-image.tar
	@echo The image is finished, Please use $(IMAGE_TAG):$(VERSION)


