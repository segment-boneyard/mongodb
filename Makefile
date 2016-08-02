
IMAGE=segment/mongodb-source
VERSION=v0.1.2-beta

save-deps:
	godep save

run:
	go run main.go mongo.go collection.go config.go description.go $(ARGS)

build:
	go build .
	
test:
	go test .

build-image:
	docker build $(FLAGS) -t $(IMAGE):$(VERSION) . 

push-image:
	docker push $(IMAGE):$(VERSION)

.PHONY: run build test build-image push-image
.DEFAULT_GOAL := build