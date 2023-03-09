compile:
	protoc api/v1/*.proto \
		--go_out=. \
		--go-grpc_out=.  \
		--go_opt=paths=source_relative \
		--go-grpc_opt=paths=source_relative \
		--proto_path=.

clean:
	rm -rf logly

build: clean compile
	go build -o logly ./cmd/main.go
