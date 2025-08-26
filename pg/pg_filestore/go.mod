module github.com/lukaszkaleta/saas-go/pg/pg_filestore

go 1.24.6

require (
	github.com/lukaszkaleta/saas-go/filestore v0.0.7
	github.com/lukaszkaleta/saas-go/pg v0.0.7
)

require (
	github.com/jackc/pgpassfile v1.0.0 // indirect
	github.com/jackc/pgservicefile v0.0.0-20240606120523-5a60cdf6a761 // indirect
	github.com/jackc/pgx/v5 v5.7.5 // indirect
	github.com/jackc/puddle/v2 v2.2.2 // indirect
	github.com/lukaszkaleta/saas-go/universal v0.0.0-20250826182527-027742bb6156 // indirect
	golang.org/x/crypto v0.37.0 // indirect
	golang.org/x/sync v0.13.0 // indirect
	golang.org/x/text v0.24.0 // indirect
)

replace github.com/lukaszkaleta/saas-go => ../../

replace github.com/lukaszkaleta/saas-go/filestore => ../../filestore

replace github.com/lukaszkaleta/saas-go/pg => ../../pg
