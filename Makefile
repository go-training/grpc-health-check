proto/health.pb.go: proto/health.proto
	protoc -I proto proto/health.proto --go_out=plugins=grpc:proto

generate_proto: proto/health.pb.go
