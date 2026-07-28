module github.com/dmi3midd/grpcsso

go 1.26.5

require (
	buf.build/go/protovalidate v1.2.0
	github.com/dmi3midd/grpcsso-protos v0.0.0-20260721122326-c830a68ff99d
	github.com/golang-jwt/jwt/v5 v5.3.1
	github.com/google/uuid v1.6.0
	github.com/jackc/pgx v3.6.2+incompatible
	github.com/jmoiron/sqlx v1.4.0
	github.com/pressly/goose/v3 v3.27.2
	github.com/redis/go-redis/v9 v9.21.0
	golang.org/x/crypto v0.54.0
	google.golang.org/grpc v1.82.0
	google.golang.org/protobuf v1.36.11
	gopkg.in/yaml.v3 v3.0.1
)

require (
	buf.build/gen/go/bufbuild/protovalidate/protocolbuffers/go v1.36.11-20260709200747-435963d16310.1 // indirect
	cel.dev/expr v0.25.1 // indirect
	github.com/antlr4-go/antlr/v4 v4.13.1 // indirect
	github.com/cespare/xxhash/v2 v2.3.0 // indirect
	github.com/cockroachdb/apd v1.1.0 // indirect
	github.com/gofrs/uuid v4.4.0+incompatible // indirect
	github.com/google/cel-go v0.28.0 // indirect
	github.com/jackc/fake v0.0.0-20150926172116-812a484cc733 // indirect
	github.com/mfridman/interpolate v0.0.2 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/sethvargo/go-retry v0.3.0 // indirect
	go.uber.org/atomic v1.11.0 // indirect
	go.uber.org/multierr v1.11.0 // indirect
	golang.org/x/exp v0.0.0-20260410095643-746e56fc9e2f // indirect
	golang.org/x/net v0.56.0 // indirect
	golang.org/x/sync v0.22.0 // indirect
	golang.org/x/sys v0.47.0 // indirect
	golang.org/x/text v0.40.0 // indirect
	google.golang.org/genproto/googleapis/api v0.0.0-20260414002931-afd174a4e478 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20260420184626-e10c466a9529 // indirect
)

replace github.com/dmi3midd/grpcsso-protos => ../grpcsso-protos
