build:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o worker *.go

docker:
	docker build -t worker:latest .
