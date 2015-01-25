# Start from a Debian image with the latest version of Go installed
# and a workspace (GOPATH) configured at /go.

FROM google/golang

WORKDIR /gopath/src/app
ADD . /gopath/src/app/
RUN go get app

CMD []
ENTRYPOINT ["/gopath/bin/app"]

# Document that the service listens on port 8080.
EXPOSE 8080
