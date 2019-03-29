ARG GOPATH=/go
ARG WORKDIR=$GOPATH/src/github.com/cakesmith/webrender

FROM znly/protoc as protoc
ARG WORKDIR
WORKDIR $WORKDIR
ADD protos protos
RUN mkdir -p protos/go protos/js
RUN protoc --go_out=./protos/go --js_out=./protos/js protos/RFB.proto

FROM golang:1.12.1-alpine3.9 as golang
ARG GOPATH
ARG WORKDIR
ENV PATH $GOPATH/bin:/usr/local/go/bin:$PATH
WORKDIR $WORKDIR
RUN apk update && apk add --no-cache git
ADD . .
COPY --from=protoc $WORKDIR/protos protos
RUN go get ./...
RUN CGO_ENABLED=0 go test ./...
RUN mkdir -p /build
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s"  -o /build/webrender

FROM alpine:3.9
ARG WORKDIR
EXPOSE $PORT
RUN apk update && apk add --no-cache ca-certificates tzdata && update-ca-certificates
COPY --from=protoc $WORKDIR/protos/js /public
COPY --from=golang $WORKDIR/public /public
COPY --from=golang build /
RUN adduser -D -g '' appuser
USER appuser
CMD ["./webrender"]