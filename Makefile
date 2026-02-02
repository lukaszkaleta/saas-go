VERSION := v0.2.176
tags:
	git add .
	git commit -m 'incremental version ...'
	git push
	git tag "universal/${VERSION}"
	git tag "database/pg/${VERSION}"
	git tag "universal/pg/${VERSION}"
	git tag "filestore/${VERSION}"
	git tag "filestore/pg/${VERSION}"
	git tag "filestore/integration/aws-s3/${VERSION}"
	git tag "messages/${VERSION}"
	git tag "messages/pg/${VERSION}"
	git tag "category/${VERSION}"
	git tag "category/pg/${VERSION}"
	git tag "job/${VERSION}"
	git tag "job/pg/${VERSION}"
	git tag "user/${VERSION}"
	git tag "user/pg/${VERSION}"
	git push --tags

reset_build:
	cd universal && go clean -modcache && rm -r go.sum || true && go mod tidy && go build && cd ..
	cd database/pg && go clean -modcache && rm -r go.sum || true && go mod tidy && go build && cd ../..
	cd universal/pg && go clean -modcache && rm -r go.sum || true && go mod tidy && go build && cd ../..
	cd filestore && go clean -modcache && rm -r go.sum || true && go mod tidy && go build && cd ..
	cd filestore/pg && go clean -modcache && rm -r go.sum || true && go mod tidy && go mod tidy && go build && cd ../..
	cd messages && go clean -modcache && rm -r go.sum || true && go mod tidy && go build && cd ..
	cd messages/pg && go clean -modcache && rm -r go.sum || true && go mod tidy && go mod tidy && go build && cd ../..
	cd category && go clean -modcache && rm -r go.sum || true && go mod tidy && go build && cd ..
	cd category/pg && go clean -modcache && rm -r go.sum || true && go mod tidy && go build && cd ../..
	cd job && go clean -modcache && rm -r go.sum || true && go mod tidy && go build && cd ..
	cd job/pg && go clean -modcache && rm -r go.sum || true && go mod tidy && go build && cd ../..
	cd user && go clean -modcache && rm -r go.sum || true && go mod tidy && go build && cd ..
	cd user/pg && go clean -modcache && rm -r go.sum || true && go mod tidy && go build && cd ../..

build:
	cd universal && go build && cd ..
	cd database/pg && go build && cd ../..
	cd universal/pg && go build && cd ../..
	cd filestore && go build && cd ..
	cd filestore/pg && go mod tidy && go build && cd ../..
	cd messages && go build && cd ..
	cd messages/pg && go mod tidy && go build && cd ../..
	cd category && go build && cd ..
	cd category/pg && go build && cd ../..
	cd job && go build && cd ..
	cd job/pg && go build && cd ../..
	cd user && go build && cd ..
	cd user/pg && go build && cd ../..

test:
	cd universal && go test && cd ..
	cd database/pg && go test && cd ../..
	cd universal/pg && go test && cd ../..
	cd filestore && go test && cd ..
	cd filestore/pg && go test && cd ../..
	cd messages && go test && cd ..
	cd messages/pg && go test && cd ../..
	cd category && go test && cd ..
	cd category/pg && go test && cd ../..
	cd job && go test && cd ..
	cd job/pg && go test && cd ../..
	cd user && go test && cd ..
	cd user/pg && go test && cd ../..

