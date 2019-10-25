FROM golang:1.10.3-alpine

# Update packages and install dependency packages for services
RUN apk update && apk add --no-cache bash git

# Change working directory
WORKDIR $GOPATH/src/eliza/

# Install dependencies
RUN go get -u github.com/golang/dep/...
RUN go get -u github.com/derekparker/delve/cmd/dlv/...
COPY . ./
RUN if test -e "Gopkg.toml"; then dep ensure -v; fi

ENV PORT 8080
ENV GIN_MODE release
EXPOSE 8080

RUN go build -o app
CMD ["./app"]