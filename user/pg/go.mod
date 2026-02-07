module github.com/lukaszkaleta/saas-go/user/pg

go 1.25.0

require (
	github.com/jackc/pgx/v5 v5.8.0
	github.com/lukaszkaleta/saas-go/database/pg v0.2.202
	github.com/lukaszkaleta/saas-go/filestore v0.2.202
	github.com/lukaszkaleta/saas-go/filestore/pg v0.2.202
	github.com/lukaszkaleta/saas-go/universal v0.2.202
	github.com/lukaszkaleta/saas-go/universal/pg v0.2.202
	github.com/lukaszkaleta/saas-go/user v0.2.202
)

require (
	github.com/jackc/pgpassfile v1.0.0 // indirect
	github.com/jackc/pgservicefile v0.0.0-20240606120523-5a60cdf6a761 // indirect
	github.com/jackc/puddle/v2 v2.2.2 // indirect
	golang.org/x/sync v0.19.0 // indirect
	golang.org/x/text v0.33.0 // indirect
)

replace github.com/lukaszkaleta/saas-go/user => ../../user

replace github.com/lukaszkaleta/saas-go/universal => ../../universal

replace github.com/lukaszkaleta/saas-go/filestore => ../../filestore

replace github.com/lukaszkaleta/saas-go/database/pg => ../../database/pg

replace github.com/lukaszkaleta/saas-go/universal/pg => ../../universal/pg

replace github.com/lukaszkaleta/saas-go/filestore/pg => ../../filestore/pg
