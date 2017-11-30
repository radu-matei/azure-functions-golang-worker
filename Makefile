plugins = grpc
target = go
protoc_location = ../azure-webjobs-sdk-script/src/WebJobs.Script.Grpc/Proto
proto_out_dir = rpc/

.PHONY: rpc
rpc:
	protoc -I $(protoc_location) --$(target)_out=plugins=$(plugins):$(proto_out_dir) $(protoc_location)/*.proto
