run:
	go run ./cmds/cap | go run ./cmds/logfmt

chat:
	go run ./cmds/chatapp | go run ./cmds/logfmt

test:
	go test ./...

tidy:
	go mod tidy
	go mod vendor
	go mod verify

upgrade:
	go get -u -v ./...
	go mod tidy
	go mod vendor
	go mod verify