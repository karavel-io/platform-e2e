FROM golang:1.17-buster as tools
RUN go install github.com/jstemmer/go-junit-report@v0.9.1

FROM golang:1.17-buster as build

COPY --from=tools /go/bin/go-junit-report /usr/local/bin/

COPY ./pkg /src/pkg
WORKDIR /src
COPY go.* /src/
ENV CGO_ENABLED=0
RUN go mod download
COPY run.sh /src
RUN --mount=type=cache,target=/root/.cache/go-build \
    go test -c -o /usr/local/bin/e2e.test ./...

CMD ["./run.sh"]
