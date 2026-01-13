module github.com/lukaszkaleta/saas-go/job/pg

go 1.25.0

require (
	github.com/jackc/pgx/v5 v5.7.6
	github.com/lukaszkaleta/saas-go/database/pg v0.2.104
	github.com/lukaszkaleta/saas-go/filestore v0.2.104
	github.com/lukaszkaleta/saas-go/filestore/pg v0.2.63
	github.com/lukaszkaleta/saas-go/job v0.2.63
	github.com/lukaszkaleta/saas-go/messages v0.2.104
	github.com/lukaszkaleta/saas-go/messages/pg v0.2.104
	github.com/lukaszkaleta/saas-go/universal v0.2.104
	github.com/lukaszkaleta/saas-go/universal/pg v0.2.104
	github.com/lukaszkaleta/saas-go/user v0.2.104
)

require (
	github.com/jackc/pgpassfile v1.0.0 // indirect
	github.com/jackc/pgservicefile v0.0.0-20240606120523-5a60cdf6a761 // indirect
	github.com/jackc/puddle/v2 v2.2.2 // indirect
	golang.org/x/crypto v0.44.0 // indirect
	golang.org/x/sync v0.18.0 // indirect
	golang.org/x/text v0.31.0 // indirect
)

replace github.com/lukaszkaleta/saas-go/filestore => ../../filestore

replace github.com/lukaszkaleta/saas-go/job => ../../job

replace github.com/lukaszkaleta/saas-go/user => ../../user

replace github.com/lukaszkaleta/saas-go/universal => ../../universal

replace github.com/lukaszkaleta/saas-go/database/pg => ../../database/pg

replace github.com/lukaszkaleta/saas-go/universal/pg => ../../universal/pg

replace github.com/lukaszkaleta/saas-go/messages => ../../messages

replace github.com/lukaszkaleta/saas-go/messages/pg => ../../messages/pg

replace github.com/lukaszkaleta/saas-go/filestore/pg => ../../filestore/pg
