FROM golang:1.18.0-buster AS builder

ARG VERSION=dev

WORKDIR /go/src/app
COPY . .
RUN go build -o raza -ldflags=-X=main.go.go.version=${VERSION} ./cmd/raza

FROM debian:buster-slim
COPY --from=builder /go/src/app/raza /go/bin/raza
ENV PATH="/go/bin:${PATH}"
CMD ["raza"]