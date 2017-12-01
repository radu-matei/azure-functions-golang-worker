plugins = grpc
target = go
protoc_location = ../azure-webjobs-sdk-script/src/WebJobs.Script.Grpc/Proto
proto_out_dir = rpc/

OUTPUT_DIR = bin
GOLANG_WORKER_BINARY = golang-worker

.PHONY: rpc
rpc:
	protoc -I $(protoc_location) --$(target)_out=plugins=$(plugins):$(proto_out_dir) $(protoc_location)/*.proto

.PHONY: golang-worker
golang-worker:
	GOOS=linux go build -o $(OUTPUT_DIR)/$(GOLANG_WORKER_BINARY)
