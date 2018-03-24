default:

check:
ifndef LOGICALID
	echo "Environment variable LOGICALID not specified."
	exit 1
endif
ifndef STACKNAME
	echo "Environment variable STACKNAME not specified."
	exit 1
endif

build:
	CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o $$GOPATH/bin/cfn-signaler github.com/marjamis/cfn-signaler

run: check
	LOGICALID=$$LOGICALID STACKNAME=$$STACKNAME go run main.go

dbuild:
	docker build -t cfn-signaler .

drun: check
ifndef PUBLICPORT
	echo "Environment variable PUBLICPORT not specified."
	exit 1
endif
	docker run -it --rm -e STACKNAME=$$STACKNAME -e LOGICALID=$$LOGICALID -p $$PUBLICPORT:8080 cfn-signaler
