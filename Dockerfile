FROM golang:1.9.2-alpine as builder
WORKDIR /go/src/github.com/salemove/node-network-detector

RUN apk add --no-cache git make && \
  go get -u github.com/Masterminds/glide/...

COPY glide.yaml glide.lock vendor ./
RUN glide install

COPY *.go ./
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o /node-network-detector .

FROM scratch
COPY --from=builder /node-network-detector /
ENTRYPOINT ["/node-network-detector"]
