FROM docker.io/curlimages/curl:latest as linkerd
ARG LINKERD_AWAIT_VERSION=v0.2.3
RUN curl -sSLo /tmp/linkerd-await https://github.com/linkerd/linkerd-await/releases/download/release%2F${LINKERD_AWAIT_VERSION}/linkerd-await-${LINKERD_AWAIT_VERSION}-amd64 && \
    chmod 755 /tmp/linkerd-await

FROM golang:1.16.3 as builder
RUN apt-get update -qq && \
    apt-get install -y -q gcc \
                        g++ \
                        make \
                        zlib1g-dev \
                        git
WORKDIR /opt
COPY go.mod .
COPY go.sum .
RUN go mod download
ADD . .
RUN go build -o clientd client/main.go
RUN go build -o serverd server/main.go

FROM ubuntu:20.04 as runner
ENV TZ 'Asia/Tehran'
COPY --from=linkerd /tmp/linkerd-await /linkerd-await
COPY --from=builder /opt/clientd /bin/clientd
COPY --from=builder /opt/serverd /bin/serverd

ENTRYPOINT ["/linkerd-await", "--shutdown", "--"]
CMD ["/bin/bash"]
