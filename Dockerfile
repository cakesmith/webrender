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
# Create appuser.
RUN adduser -D -g '' appuser
RUN apk update && apk add --no-cache git ca-certificates tzdata && update-ca-certificates
ADD . .
COPY --from=protoc $WORKDIR/protos protos
RUN go get ./...
RUN CGO_ENABLED=0 go test ./...
RUN mkdir -p /build
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s"  -o /build/webrender

#FROM alpine:3.9
FROM scratch
ARG WORKDIR
EXPOSE $PORT
#ENV PORT=$PORT
COPY --from=protoc $WORKDIR/protos/js /public
COPY --from=golang /usr/share/zoneinfo /usr/share/zoneinfo
COPY --from=golang /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=golang /etc/passwd /etc/passwd
COPY --from=golang $WORKDIR/public /public
COPY --from=golang build /
#RUN useradd -ms /bin/bash myuser
USER appuser
CMD ["./webrender"]