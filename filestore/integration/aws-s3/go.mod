module github.com/lukaszkaleta/saas-go/filestore/integration/aws-s3

go 1.24.6

require (
	github.com/lukaszkaleta/saas-go/filestore v0.1.8

	github.com/aws/aws-sdk-go-v2 v1.38.1
    	github.com/aws/aws-sdk-go-v2/config v1.31.3
    	github.com/aws/aws-sdk-go-v2/service/s3 v1.87.1
)

replace github.com/lukaszkaleta/saas-go/filestore => ../../../filestore
