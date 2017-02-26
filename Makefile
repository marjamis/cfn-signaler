build:
	CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o ./bin/cfn-signaler.go docker_cfn-signaler_go
	docker build -t cfn-signaler .
