plugins = grpc
target = go
protoc_location = rpc/
proto_out_dir = rpc/

OUTPUT_DIR = sample
GOLANG_WORKER_BINARY = golang-worker
SUBDIRS := $(wildcard sample/*)

.PHONY: rpc
rpc:
        protoc -I $(protoc_location) --$(target)_out=plugins=$(plugins):$(proto_out_dir) $(protoc_location)/*.proto

.PHONY: golang-worker
golang-worker:
        GOOS=linux go build -o $(OUTPUT_DIR)/$(GOLANG_WORKER_BINARY)

.PHONY: dep
dep:
	go get -u github.com/golang/dep/... && \
	dep ensure

.PHONY : samples $(SUBDIRS)
samples : $(SUBDIRS)

$(SUBDIRS) :
	cd $@ && \
	go build -buildmode=plugin
