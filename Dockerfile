ARG GOPATH=/go
ARG WORKDIR=$GOPATH/src/github.com/cakesmith/webrender

FROM znly/protoc as protoc
ARG WORKDIR
WORKDIR $WORKDIR
ADD protos protos
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
RUN mkdir /build
RUN CGO_ENABLED=0 go build -o /build/webrender

FROM heroku/heroku:18
ARG WORKDIR
EXPOSE $PORT
ENV PORT=$PORT
COPY --from=protoc $WORKDIR/protos/js ./public
COPY --from=golang $WORKDIR/public ./public
COPY --from=golang /build .
CMD ["./webrender"]