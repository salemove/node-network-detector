FROM golang:1.9.2 as builder
WORKDIR /go/src/github.com/salemove/node-network-detector
COPY pinger.go .
RUN go get -d -v
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o /node-network-detector .

FROM scratch
COPY --from=builder /node-network-detector /
ENTRYPOINT ["/node-network-detector"]
