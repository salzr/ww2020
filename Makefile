IMAGE=salzr/cloudrun-storage
VERS_TAG=v0.0.0

build:
	docker build --build-arg VERS_TAG=$(VERS_TAG) -t $(IMAGE):$(VERS_TAG) .