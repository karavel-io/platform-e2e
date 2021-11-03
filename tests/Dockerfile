FROM golang:1.17-buster as tools
RUN go install github.com/jstemmer/go-junit-report@v0.9.1

FROM golang:1.17-buster as build

COPY --from=tools /go/bin/go-junit-report /usr/local/bin/

# Install kubectl
# Note: Latest version may be found on:
# https://aur.archlinux.org/packages/kubectl-bin/
RUN wget https://storage.googleapis.com/kubernetes-release/release/v1.21.3/bin/linux/amd64/kubectl -O /usr/bin/kubectl && \
    chmod +x /usr/bin/kubectl && \
    apt-get update && \
    apt-get install -y jq

COPY ./pkg /src/pkg
WORKDIR /src
COPY go.* /src/
ENV CGO_ENABLED=0
RUN go mod download
COPY run.sh /src
RUN --mount=type=cache,target=/root/.cache/go-build \
    go test -c -o /usr/local/bin/e2e.test ./...

CMD ["./run.sh"]
