module github.com/lukaszkaleta/saas-go/pg/pg_filestore

go 1.24.6

require (
	github.com/lukaszkaleta/saas-go/filestore v.0.0.7
	github.com/lukaszkaleta/saas-go/pg v.0.0.7
)
replace github.com/lukaszkaleta/saas-go => ../../
replace github.com/lukaszkaleta/saas-go/filestore => ../../filestore
replace github.com/lukaszkaleta/saas-go/pg => ../../pg