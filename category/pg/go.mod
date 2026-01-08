module github.com/lukaszkaleta/saas-go/category/pg

go 1.25.0

require (
	github.com/jackc/pgx/v5 v5.7.5
	github.com/lukaszkaleta/saas-go/category v0.1.43
	github.com/lukaszkaleta/saas-go/database/pg v0.2.97
	github.com/lukaszkaleta/saas-go/universal v0.2.97
	github.com/lukaszkaleta/saas-go/universal/pg v0.1.43
)

require (
	github.com/jackc/pgpassfile v1.0.0 // indirect
	github.com/jackc/pgservicefile v0.0.0-20240606120523-5a60cdf6a761 // indirect
	github.com/jackc/puddle/v2 v2.2.2 // indirect
	golang.org/x/crypto v0.41.0 // indirect
	golang.org/x/sync v0.16.0 // indirect
	golang.org/x/text v0.28.0 // indirect
)

replace github.com/lukaszkaleta/saas-go/category => ../../category

replace github.com/lukaszkaleta/saas-go/universal => ../../universal

replace github.com/lukaszkaleta/saas-go/database/pg => ../../database/pg

replace github.com/lukaszkaleta/saas-go/universal/pg => ../../universal/pg
