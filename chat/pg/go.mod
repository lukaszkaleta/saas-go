module github.com/lukaszkaleta/saas-go/chat/pg

go 1.25.0

require (
	github.com/jackc/pgx/v5 v5.9.1
	github.com/lukaszkaleta/saas-go/database/pg v0.2.315
	github.com/lukaszkaleta/saas-go/chat v0.2.315
	github.com/lukaszkaleta/saas-go/universal v0.2.315
)

replace github.com/lukaszkaleta/saas-go/database/pg => ../../database/pg
replace github.com/lukaszkaleta/saas-go/universal => ../../universal
replace github.com/lukaszkaleta/saas-go/chat => ../../chat
