module github.com/tjamir/flisol-2025/microblog/timeline-worker

go 1.24.1

require github.com/gocql/gocql v1.7.0

require (
	github.com/cespare/xxhash/v2 v2.3.0 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/dgryski/go-rendezvous v0.0.0-20200823014737-9f7001d12a5f // indirect
	github.com/klauspost/compress v1.15.9 // indirect
	github.com/pierrec/lz4/v4 v4.1.15 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	golang.org/x/net v0.35.0 // indirect
	golang.org/x/sys v0.30.0 // indirect
	golang.org/x/text v0.22.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20250422160041-2d3770c4ea7f // indirect
	google.golang.org/protobuf v1.36.6 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

require (
	github.com/golang/snappy v0.0.3 // indirect
	github.com/hailocab/go-hostpool v0.0.0-20160125115350-e80d13ce29ed // indirect
	github.com/redis/go-redis/v9 v9.7.3
	github.com/segmentio/kafka-go v0.4.47
	github.com/stretchr/testify v1.10.0
	github.com/tjamir/flisol-2025/microblog/follow-service v0.0.0-00010101000000-000000000000
	github.com/tjamir/flisol-2025/microblog/post-service v0.0.0-00010101000000-000000000000
	google.golang.org/grpc v1.72.0
	gopkg.in/inf.v0 v0.9.1 // indirect
)

replace github.com/tjamir/flisol-2025/microblog/post-service => ../post-service

replace github.com/tjamir/flisol-2025/microblog/follow-service => ../follow-service
