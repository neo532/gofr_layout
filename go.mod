module github.com/neo532/gofr_layout

go 1.26.1

require (
	github.com/google/wire v0.5.0
	github.com/neo532/gofr v1.0.16-0.20260620132837-165450a018e1
	github.com/neo532/gofr/transport/http v0.0.0-20260620132837-165450a018e1
	github.com/neo532/gofr_layout/proto v0.0.0
	github.com/neo532/gokit v1.0.45
	github.com/neo532/gokit/database/orm v1.0.0
	github.com/neo532/gokit/database/redis v1.0.0
	github.com/neo532/gokit/filepath v1.0.0
	github.com/neo532/gokit/logger/writer/lumberjack v1.0.0
	gopkg.in/yaml.v3 v3.0.1
	gorm.io/driver/mysql v1.6.0
)

require (
	filippo.io/edwards25519 v1.1.0 // indirect
	github.com/cespare/xxhash/v2 v2.3.0 // indirect
	github.com/dgryski/go-rendezvous v0.0.0-20200823014737-9f7001d12a5f // indirect
	github.com/envoyproxy/protoc-gen-validate v1.3.3 // indirect
	github.com/fsnotify/fsnotify v1.10.1 // indirect
	github.com/go-redis/redis/v8 v8.11.5 // indirect
	github.com/go-sql-driver/mysql v1.8.1 // indirect
	github.com/google/subcommands v1.0.1 // indirect
	github.com/jinzhu/inflection v1.0.0 // indirect
	github.com/jinzhu/now v1.1.5 // indirect
	github.com/julienschmidt/httprouter v1.3.0 // indirect
	github.com/neo532/gofr/transport/grpc v0.0.0-20260620132837-165450a018e1 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	golang.org/x/mod v0.32.0 // indirect
	golang.org/x/net v0.51.0 // indirect
	golang.org/x/sync v0.20.0 // indirect
	golang.org/x/sys v0.42.0 // indirect
	golang.org/x/text v0.34.0 // indirect
	golang.org/x/tools v0.41.0 // indirect
	google.golang.org/genproto/googleapis/api v0.0.0-20260618152121-87f3d3e198d3 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20260610212136-7ab31c22f7ad // indirect
	google.golang.org/grpc v1.81.1 // indirect
	google.golang.org/protobuf v1.36.11 // indirect
	gopkg.in/natefinch/lumberjack.v2 v2.2.1 // indirect
	gorm.io/gorm v1.31.1 // indirect
)

replace (
	github.com/neo532/gofr_layout/proto => ./proto
	github.com/neo532/gokit/database/orm => ../gokit/database/orm
	github.com/neo532/gokit/database/redis => ../gokit/database/redis
)
