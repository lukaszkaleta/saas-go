module github.com/lukaszkaleta/saas-go/pg/pg_filestore

go 1.24.6

require (
	github.com/lukaszkaleta/saas-go/filestore v0.0.5
	github.com/lukaszkaleta/saas-go/pg v0.0.5
)
replace github.com/lukaszkaleta/saas-go => ../../
replace github.com/lukaszkaleta/saas-go/filestore => ../../filestore
replace github.com/lukaszkaleta/saas-go/pg => ../../pg