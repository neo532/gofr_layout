module github.com/neo532/gofr_layout/proto

go 1.25.0

require (
	github.com/envoyproxy/protoc-gen-validate v1.3.3
	github.com/neo532/gofr v0.0.0
	github.com/neo532/gofr/transport/grpc v0.0.0-00010101000000-000000000000
	github.com/neo532/gofr/transport/http v0.0.0-00010101000000-000000000000
	google.golang.org/genproto/googleapis/api v0.0.0-20260618152121-87f3d3e198d3
	google.golang.org/protobuf v1.36.11
)

require (
	github.com/julienschmidt/httprouter v1.3.0 // indirect
	golang.org/x/net v0.51.0 // indirect
	golang.org/x/sys v0.42.0 // indirect
	golang.org/x/text v0.34.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20260610212136-7ab31c22f7ad // indirect
	google.golang.org/grpc v1.81.1 // indirect
)

replace (
	github.com/neo532/gofr => ../../gofr
	github.com/neo532/gofr/transport/grpc => ../../gofr/transport/grpc
	github.com/neo532/gofr/transport/http => ../../gofr/transport/http
)
