FROM scratch

COPY ./application/bin/cfn-signaler.go /entrypoint
COPY ./application/src/cfn-signaler/stylesheets/* /stylesheets/
COPY ./application/src/cfn-signaler/templates/* /templates/

ENTRYPOINT [ "/entrypoint" ]
