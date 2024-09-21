
run-server:
	go run ./cmd/server/main.go

run-client:
	go run ./cmd/client/main.go

run-resource:
	go run ./cmd/resource/main.go

test:
	go test ./...

build-server:
	go build  -o ./build/server/server ./cmd/server/main.go

up-server:
	./build/server/server