FROM golang:1.9.2-alpine as builder
WORKDIR /go/src/github.com/salemove/node-network-detector

COPY glide.yaml glide.lock ./
RUN apk add --update --no-cache wget curl git \
    && wget "https://github.com/Masterminds/glide/releases/download/v0.12.3/glide-v0.12.3-`go env GOHOSTOS`-`go env GOHOSTARCH`.tar.gz" -O /tmp/glide.tar.gz \
    && mkdir /tmp/glide \
    && tar --directory=/tmp/glide -xvf /tmp/glide.tar.gz \
    && rm -rf /tmp/glide.tar.gz \
    && export PATH=$PATH:/tmp/glide/`go env GOHOSTOS`-`go env GOHOSTARCH` \
    && glide update -v \
    && glide install

COPY pinger.go .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o /node-network-detector .

FROM scratch
COPY --from=builder /node-network-detector /
ENTRYPOINT ["/node-network-detector"]
