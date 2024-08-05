module github.com/mcq_backend

go 1.22.5

require (
	github.com/bkend-db v0.0.1
	github.com/bkend-redis v0.0.1
	github.com/google/uuid v1.6.0
	github.com/gorilla/mux v1.8.1
	github.com/joho/godotenv v1.5.1
	github.com/nats-io/nats.go v1.33.1
)

replace github.com/bkend-redis => ../bkend-redis

replace github.com/bkend-db => ../bkend-db

require (
	github.com/klauspost/compress v1.17.6 // indirect
	github.com/nats-io/nkeys v0.4.7 // indirect
	github.com/nats-io/nuid v1.0.1 // indirect
	golang.org/x/crypto v0.22.0 // indirect
	golang.org/x/sys v0.19.0 // indirect
)
