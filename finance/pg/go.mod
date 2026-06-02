module github.com/lukaszkaleta/saas-go/finance/pg

go 1.24.0

require (
	github.com/lukaszkaleta/saas-go/database/pg v0.2.313
	github.com/lukaszkaleta/saas-go/finance v0.2.313
	github.com/lukaszkaleta/saas-go/universal v0.2.313
)

replace github.com/lukaszkaleta/saas-go/database/pg => ../../database/pg
replace github.com/lukaszkaleta/saas-go/universal => ../../universal
replace github.com/lukaszkaleta/saas-go/finance => ../../finance
