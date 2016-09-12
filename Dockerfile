FROM alpine:3.3

RUN apk --update upgrade && apk add ca-certificates && update-ca-certificates

COPY ./application/bin/cfn-signaler.go /entrypoint
COPY ./application/src/cfn-signaler/stylesheets/* /stylesheets/
COPY ./application/src/cfn-signaler/templates/* /templates/

WORKDIR /

ENTRYPOINT [ "/entrypoint" ]
