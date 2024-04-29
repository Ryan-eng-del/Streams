
all: generate

generate: test/test.proto
	  protoc \
	  --go-grpc_out=. \
		--go_out=. \
		test/test.proto

