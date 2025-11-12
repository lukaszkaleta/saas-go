VERSION := v0.2.68
tags:
	git add .
	git commit -m 'incremental version ...'
	git push
	git tag "database/pg/${VERSION}"
	git tag "universal/${VERSION}"
	git tag "universal/pg/${VERSION}"
	git tag "filestore/${VERSION}"
	git tag "filestore/pg/${VERSION}"
	git tag "filestore/integration/aws-s3/${VERSION}"
	git tag "category/${VERSION}"
	git tag "category/pg/${VERSION}"
	git tag "job/${VERSION}"
	git tag "job/pg/${VERSION}"
	git tag "user/${VERSION}"
	git tag "user/pg/${VERSION}"
	git push --tags

build:
	cd universal && go mod tidy && go build && cd ..
	cd database/pg && go mod tidy && go build && cd ../..
	cd universal/pg && go mod tidy && go build && cd ../..
	cd filestore && go mod tidy && go build && cd ..
	cd filestore/pg && go mod tidy && go mod tidy && go build && cd ../..
	cd category && go mod tidy && go build && cd ..
	cd category/pg && go mod tidy && go build && cd ../..
	cd job && go mod tidy && go build && cd ..
	cd job/pg && go mod tidy && go build && cd ../..
	cd user && go mod tidy && go build && cd ..
	cd user/pg && go mod tidy && go build && cd ../..

