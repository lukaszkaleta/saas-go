module github.com/lukaszkaleta/saas-go/pg/pg-universal

go 1.24.6

require (
	github.com/lukaszkaleta/saas-go/universal v0.1.0
	github.com/lukaszkaleta/saas-go/pg/database v0.1.1
)

replace github.com/lukaszkaleta/saas-go/universal => ../../universal
replace github.com/lukaszkaleta/saas-go/pg/database => ../database
