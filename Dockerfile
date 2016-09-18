FROM alpine:3.3

RUN apk --no-cache --update upgrade && apk add --no-cache ca-certificates && update-ca-certificates

COPY ./application/bin/cfn-signaler.go /entrypoint
COPY ./application/src/cfn-signaler/stylesheets/* /stylesheets/
COPY ./application/src/cfn-signaler/templates/* /templates/

WORKDIR /

ENTRYPOINT [ "/entrypoint" ]
