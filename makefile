run:
	go run cmds/cap/main.go

test:
	go test ./...

tidy:
	go mod tidy
	go mod vendor

upgrade:
	go get -u -v ./...
	go mod tidy
	go mod vendor