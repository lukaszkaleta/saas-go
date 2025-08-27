module github.com/lukaszkaleta/saas-go/universal/pg

go 1.24.6

require (
	github.com/lukaszkaleta/saas-go/universal v0.1.6
	github.com/lukaszkaleta/saas-go/database/pg v0.1.6
)

replace github.com/lukaszkaleta/saas-go/universal => ../../universal
replace github.com/lukaszkaleta/saas-go/database/pg => ../../database/pg
