FROM alpine:3.3

RUN apk --no-cache --update upgrade && apk add --no-cache ca-certificates && update-ca-certificates

COPY ./bin/cfn-signaler.go /entrypoint
COPY ./stylesheets/* /stylesheets/
COPY ./templates/* /templates/

WORKDIR /

ENTRYPOINT [ "/entrypoint" ]
