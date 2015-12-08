# Start from a Debian image with the latest version of Go installed
# and a workspace (GOPATH) configured at /go.
FROM golang

# Copy the local package files to the container's workspace.
ADD . $GOPATH/src/github.com/EconomistDigitalSolutions/goberry

WORKDIR /$GOPATH/src/github.com/EconomistDigitalSolutions/goberry

ENV GO15VENDOREXPERIMENT 1 

# Build the outyet command inside the container.
RUN go get -u github.com/lobatt/go-junit-report && \
    go get -u github.com/golang/lint/golint && \
    go get -u github.com/GeertJohan/fgt && \
    go get golang.org/x/tools/cmd/vet && \
    go install

# Run the outyet command by default when the container starts.

COPY service.conf /etc/supervisor/conf.d/

CMD ["/usr/bin/supervisord", "-c", "/etc/supervisor/conf.d/supervisord.conf"]

# Document that the service listens on port 9494.
EXPOSE 9494
