module github.com/lukaszkaleta/saas-go/user/pg

go 1.24.6

require (
	github.com/lukaszkaleta/saas-go/database/pg v0.1.8
	github.com/lukaszkaleta/saas-go/universal/pg v0.1.8
	github.com/lukaszkaleta/saas-go/user v0.1.8
)

require (
	github.com/jackc/pgpassfile v1.0.0 // indirect
	github.com/jackc/pgservicefile v0.0.0-20240606120523-5a60cdf6a761 // indirect
	github.com/jackc/pgx/v5 v5.7.5 // indirect
	github.com/jackc/puddle/v2 v2.2.2 // indirect
	github.com/lukaszkaleta/saas-go/universal v0.1.8 // indirect
	golang.org/x/crypto v0.41.0 // indirect
	golang.org/x/sync v0.16.0 // indirect
	golang.org/x/text v0.28.0 // indirect
)

replace github.com/lukaszkaleta/saas-go/user => ../../user

replace github.com/lukaszkaleta/saas-go/database/pg => ../../database/pg

replace github.com/lukaszkaleta/saas-go/universal/pg => ../../universal/pg
