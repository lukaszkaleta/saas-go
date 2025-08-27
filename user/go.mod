module github.com/lukaszkaleta/saas-go/user

go 1.24.6

require github.com/lukaszkaleta/saas-go/universal v0.1.7
require github.com/lukaszkaleta/saas-go/filestore v0.1.7

replace github.com/lukaszkaleta/saas-go => ../

replace github.com/lukaszkaleta/saas-go/universal => ../universal
replace github.com/lukaszkaleta/saas-go/filestore => ../filestore
