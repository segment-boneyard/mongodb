
IMAGE=segment/mongodb-source
VERSION=v0.1.4-beta

save-deps:
	godep save

run:
	godep go run main.go $(ARGS)

build:
	godep go build .
	
test:
	godep go test .

build-image:
	docker build $(FLAGS) -t $(IMAGE):$(VERSION) . 

push-image:
	docker push $(IMAGE):$(VERSION)

.PHONY: run build test build-image push-image
.DEFAULT_GOAL := build