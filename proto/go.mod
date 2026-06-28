module github.com/neo532/gofr_layout/proto

go 1.26.0

require (
	github.com/envoyproxy/protoc-gen-validate v1.3.3
	github.com/gorilla/websocket v1.5.3
	github.com/neo532/gofr v0.0.0
	github.com/neo532/gofr/transport/grpc v0.0.0-00010101000000-000000000000
	github.com/neo532/gofr/transport/http v0.0.0-00010101000000-000000000000
	github.com/neo532/gofr/transport/rpcx v0.0.0-00010101000000-000000000000
	github.com/neo532/gofr/transport/websocket v0.0.0-00010101000000-000000000000
	github.com/smallnest/rpcx v1.9.4
	google.golang.org/genproto/googleapis/api v0.0.0-20260618152121-87f3d3e198d3
	google.golang.org/grpc v1.81.1
	google.golang.org/protobuf v1.36.11
)

require (
	github.com/akutz/memconn v0.1.0 // indirect
	github.com/alitto/pond v1.9.2 // indirect
	github.com/apache/thrift v0.23.0 // indirect
	github.com/cenk/backoff v2.2.1+incompatible // indirect
	github.com/cenkalti/backoff v2.2.1+incompatible // indirect
	github.com/dgryski/go-jump v0.0.0-20211018200510-ba001c3ffce0 // indirect
	github.com/edwingeng/doublejump v1.0.1 // indirect
	github.com/facebookgo/clock v0.0.0-20150410010913-600d898af40a // indirect
	github.com/fatih/color v1.18.0 // indirect
	github.com/go-ping/ping v1.2.0 // indirect
	github.com/godzie44/go-uring v0.0.0-20220926161041-69611e8b13d5 // indirect
	github.com/gogo/protobuf v1.3.2 // indirect
	github.com/golang/snappy v0.0.4 // indirect
	github.com/google/uuid v1.6.0 // indirect
	github.com/grandcat/zeroconf v1.0.0 // indirect
	github.com/hashicorp/errwrap v1.1.0 // indirect
	github.com/hashicorp/go-multierror v1.1.1 // indirect
	github.com/hashicorp/golang-lru v1.0.2 // indirect
	github.com/juju/ratelimit v1.0.2 // indirect
	github.com/julienschmidt/httprouter v1.3.0 // indirect
	github.com/kavu/go_reuseport v1.5.0 // indirect
	github.com/klauspost/cpuid/v2 v2.2.9 // indirect
	github.com/klauspost/reedsolomon v1.12.4 // indirect
	github.com/libp2p/go-sockaddr v0.2.0 // indirect
	github.com/mattn/go-colorable v0.1.14 // indirect
	github.com/mattn/go-isatty v0.0.20 // indirect
	github.com/miekg/dns v1.1.63 // indirect
	github.com/philhofer/fwd v1.1.3-0.20240916144458-20a13a1f6b7c // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/quic-go/quic-go v0.59.1 // indirect
	github.com/rs/cors v1.11.1 // indirect
	github.com/rubyist/circuitbreaker v2.2.1+incompatible // indirect
	github.com/smallnest/gordma v0.3.0 // indirect
	github.com/smallnest/quick v0.2.0 // indirect
	github.com/soheilhy/cmux v0.1.5 // indirect
	github.com/templexxx/cpufeat v0.0.0-20180724012125-cef66df7f161 // indirect
	github.com/templexxx/xor v0.0.0-20191217153810-f85b25db303b // indirect
	github.com/tinylib/msgp v1.2.5 // indirect
	github.com/tjfoc/gmsm v1.4.1 // indirect
	github.com/valyala/fastrand v1.1.0 // indirect
	github.com/vmihailenco/msgpack/v5 v5.4.1 // indirect
	github.com/vmihailenco/tagparser/v2 v2.0.0 // indirect
	github.com/xtaci/kcp-go v5.4.20+incompatible // indirect
	golang.org/x/crypto v0.50.0 // indirect
	golang.org/x/mod v0.34.0 // indirect
	golang.org/x/net v0.53.0 // indirect
	golang.org/x/sync v0.20.0 // indirect
	golang.org/x/sys v0.43.0 // indirect
	golang.org/x/text v0.36.0 // indirect
	golang.org/x/tools v0.43.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20260610212136-7ab31c22f7ad // indirect
)

replace (
	github.com/neo532/gofr => ../../gofr
	github.com/neo532/gofr/transport/grpc => ../../gofr/transport/grpc
	github.com/neo532/gofr/transport/http => ../../gofr/transport/http
	github.com/neo532/gofr/transport/rpcx => ../../gofr/transport/rpcx
	github.com/neo532/gofr/transport/websocket => ../../gofr/transport/websocket
)
