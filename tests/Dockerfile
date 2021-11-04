FROM golang:1.17-buster as tools

FROM golang:1.17-buster as build

COPY ./pkg /src/pkg
WORKDIR /src
COPY go.* /src/
ENV CGO_ENABLED=0
RUN go mod download
COPY run.sh /src
RUN --mount=type=cache,target=/root/.cache/go-build \
    go test -c -o /usr/local/bin/e2e.test ./...

CMD ["./run.sh"]
