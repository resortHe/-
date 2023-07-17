gen:
	protoc --proto_path=proto --go_out=plugins=grpc:./proto ./proto/*.proto


clean:
	rm -rf proto/*.go

server:
	go run cmd/server/main.go -port 8080

client:
	go run cmd/client/main.go -address 0.0.0.0:8080

test:
	go test -cover -race ./...