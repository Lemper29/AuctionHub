module github.com/Lemper29/api-gateway

go 1.24.3

require (
	github.com/Lemper29/auction v0.0.0
	github.com/gorilla/mux v1.8.1
	github.com/grpc-ecosystem/grpc-gateway/v2 v2.27.2
	google.golang.org/grpc v1.75.1
)

require google.golang.org/genproto/googleapis/api v0.0.0-20250908214217-97024824d090 // indirect

replace github.com/Lemper29/auction => ../

require (
	github.com/joho/godotenv v1.5.1
	golang.org/x/net v0.44.0 // indirect
	golang.org/x/sys v0.36.0 // indirect
	golang.org/x/text v0.29.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20250826171959-ef028d996bc1 // indirect
	google.golang.org/protobuf v1.36.9 // indirect
)
