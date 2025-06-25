# Protoc-Gen-Qkrpc
A protoc plugin to generate client stubs from proto for qkrpc framework

## Installation

Make sure Go 1.20+ is installed and in your PATH.

```bash
go install github.com/ajaypanthagani/protoc-gen-qkrpc@HEAD
```

## Generate stubs

```bash
protoc \
  --go_out=. \
  --go_opt=paths=source_relative \
  --qkrpc_out=. \
  --qkrpc_opt=paths=source_relative \
  example/proto/echo.proto
```