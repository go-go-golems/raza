//go:build tools
// +build tools

package tools

import (
	_ "github.com/pressly/goose/v3/cmd/goose"
	_ "google.golang.org/grpc/cmd/protoc-gen-go-grpc"
	_ "google.golang.org/protobuf/cmd/protoc-gen-go"
)

//go:generate ./generate.sh
