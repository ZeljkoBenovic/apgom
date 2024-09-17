run:
	go run main.go -ami-host 10.255.129.9 -ami-user cxpanel -ami-pass cxmanager*con

build:
	CGO_ENABLED=0 GOOS="linux" GOARCH="amd64" go build -o apgom -ldflags '-s -w -extldflags "-static"' main.go
