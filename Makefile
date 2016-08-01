
IMAGE=segment/mongodb-source

save-deps:
	godep save

run:
	go run main.go mongo.go collection.go config.go description.go

build:
	go build .
	
test:
	go test .

build-image:
	docker build -t $(IMAGE) . 

push-image:
	docker push $(IMAGE)

.PHONY: run build test build-image push-image
.DEFAULT_GOAL := build
