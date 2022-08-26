FROM golang:1.13
MAINTAINER Gursimran Singh <singhgursimran@me.com>

# Set go bin which doesn't appear to be set already.
ENV GOBIN /go/bin
ARG BUILD_ID
ENV BUILD_IMAGE=$BUILD_ID

# build directories
ADD . /go/src/git.xenonstack.com/xs-onboarding/document-manage
WORKDIR /go/src/git.xenonstack.com/xs-onboarding/document-manage

# Go dep!
#RUN go get -u github.com/golang/dep/...
#RUN dep ensure -update

RUN go install git.xenonstack.com/xs-onboarding/document-manage

EXPOSE 8001

