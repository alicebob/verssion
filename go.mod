module github.com/alicebob/verssion

go 1.18

require (
	github.com/alicebob/minipg v0.0.0-20220802091335-ff22a5aead93
	github.com/google/uuid v1.1.2
	github.com/jackc/pgx/v4 v4.14.1
	github.com/julienschmidt/httprouter v1.3.0
	golang.org/x/net v0.0.0-20220621193019-9d032be2e588
)

require (
	github.com/jackc/chunkreader/v2 v2.0.1 // indirect
	github.com/jackc/pgconn v1.10.1 // indirect
	github.com/jackc/pgio v1.0.0 // indirect
	github.com/jackc/pgpassfile v1.0.0 // indirect
	github.com/jackc/pgproto3/v2 v2.2.0 // indirect
	github.com/jackc/pgservicefile v0.0.0-20200714003250-2b9c44734f2b // indirect
	github.com/jackc/pgtype v1.9.1 // indirect
	github.com/jackc/puddle v1.2.0 // indirect
	github.com/lib/pq v1.10.4 // indirect
	golang.org/x/crypto v0.0.0-20210711020723-a769d52b0f97 // indirect
	golang.org/x/text v0.3.7 // indirect
)

replace github.com/alicebob/minipg@main => ../minipg
