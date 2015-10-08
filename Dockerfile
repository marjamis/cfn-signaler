FROM centos:centos6

MAINTAINER J M <marjamis@amazon.com>

RUN yum update -y && yum install -y tar git 
RUN curl https://storage.googleapis.com/golang/go1.5.1.linux-amd64.tar.gz -o /tmp/go1.5.1.linux-amd64.tar.gz && \
  cd /tmp && \
  tar -C /usr/local -xzf go1.5.1.linux-amd64.tar.gz

ENV PATH $PATH:/usr/local/go/bin
ENV GOPATH /app

RUN mkdir /app && cd /app && \
  /usr/local/go/bin/go get github.com/aws/aws-sdk-go/service/cloudformation && \
  /usr/local/go/bin/go get github.com/aws/aws-sdk-go/aws/ec2metadata &&\
  /usr/local/go/bin/go get encoding/json

ADD files/main.go /app/
ADD files/stylesheets/* /app/stylesheets/
ADD files/templates/* /app/templates/

WORKDIR /app/

ENTRYPOINT [ "/usr/local/go/bin/go", "run", "./main.go" ]
