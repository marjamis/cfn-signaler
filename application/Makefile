
clean:
	rm ./bin/cfn-signaler.go

test:
	cd src/cfn-signaler/; PORT=8080 go run main.go

setup:
	GOPATH=$(shell pwd) go get "github.com/aws/aws-sdk-go/aws" "github.com/aws/aws-sdk-go/service/cloudformation" "github.com/aws/aws-sdk-go/aws/ec2metadata" "github.com/aws/aws-sdk-go/aws/session" "github.com/Sirupsen/logrus"

build:
	GOPATH=$(shell pwd) CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o ./bin/cfn-signaler.go cfn-signaler
