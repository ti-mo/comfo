package main

//go:generate protoc rpc/comfo/service.proto --twirp_out=paths=source_relative:. --go_out=paths=source_relative:.
