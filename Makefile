run:
	go run main.go -ami-host 172.28.1.10 -ami-user apgom -ami-pass 7d86e39a4d67f55969847c2ec4c07fb4

build:
	CGO_ENABLED=0 GOOS="linux" GOARCH="amd64" go build -o apgom -ldflags '-s -w -extldflags "-static"' main.go
