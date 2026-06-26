module github.com/dmi3midd/grpcsso

go 1.26.4

require (
	github.com/dmi3midd/grpcsso-protos v0.0.0-20260624162619-42be31661225
	github.com/golang-jwt/jwt/v5 v5.3.1
	github.com/jackc/pgx v3.6.2+incompatible
	github.com/jmoiron/sqlx v1.4.0
	github.com/pressly/goose/v3 v3.27.1
	github.com/redis/go-redis/v9 v9.20.1
	google.golang.org/grpc v1.81.1
	gopkg.in/yaml.v3 v3.0.1
)

require (
	github.com/cespare/xxhash/v2 v2.3.0 // indirect
	github.com/cockroachdb/apd v1.1.0 // indirect
	github.com/gofrs/uuid v4.4.0+incompatible // indirect
	github.com/jackc/fake v0.0.0-20150926172116-812a484cc733 // indirect
	github.com/mfridman/interpolate v0.0.2 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/rs/xid v1.6.0
	github.com/sethvargo/go-retry v0.3.0 // indirect
	go.uber.org/atomic v1.11.0 // indirect
	go.uber.org/multierr v1.11.0 // indirect
	golang.org/x/crypto v0.53.0
	golang.org/x/net v0.55.0 // indirect
	golang.org/x/sync v0.21.0 // indirect
	golang.org/x/sys v0.46.0 // indirect
	golang.org/x/text v0.38.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20260420184626-e10c466a9529 // indirect
	google.golang.org/protobuf v1.36.11 // indirect
)

replace github.com/dmi3midd/grpcsso-protos => ../grpcsso-protos
