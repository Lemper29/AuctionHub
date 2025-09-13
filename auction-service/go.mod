module github/auctiongithub/auction-service

go 1.24.3

require (
	github.com/google/uuid v1.6.0
	github.com/joho/godotenv v1.5.1
	google.golang.org/grpc v1.75.1
	google.golang.org/protobuf v1.36.9 // indirect
	gorm.io/driver/postgres v1.6.0
	gorm.io/gorm v1.30.5
)

require (
	github.com/jackc/pgpassfile v1.0.0 // indirect
	github.com/jackc/pgservicefile v0.0.0-20240606120523-5a60cdf6a761 // indirect
	github.com/jackc/pgx/v5 v5.7.6 // indirect
	github.com/jackc/puddle/v2 v2.2.2 // indirect
	github.com/jinzhu/inflection v1.0.0 // indirect
	github.com/jinzhu/now v1.1.5 // indirect
	github/auctiongithub/proto v0.0.0-00010101000000-000000000000
	golang.org/x/crypto v0.42.0 // indirect
	golang.org/x/net v0.44.0 // indirect
	golang.org/x/sync v0.17.0 // indirect
	golang.org/x/sys v0.36.0 // indirect
	golang.org/x/text v0.29.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20250908214217-97024824d090 // indirect
)

replace github/auctiongithub/proto => ../proto
