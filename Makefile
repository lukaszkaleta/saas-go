VERSION := v0.2.93
tags:
	git add .
	git commit -m 'incremental version ...'
	git push
	git tag "saas-${VERSION}"
	git push --tags

reset_build:
	cd universal && rm -r go.sum || true && go mod tidy && go build && cd ..
	cd database/pg && rm -r go.sum || true && go mod tidy && go build && cd ../..
	cd universal/pg && rm -r go.sum || true && go mod tidy && go build && cd ../..
	cd filestore && rm -r go.sum || true && go mod tidy && go build && cd ..
	cd filestore/pg && rm -r go.sum || true && go mod tidy && go mod tidy && go build && cd ../..
	cd messages && rm -r go.sum || true && go mod tidy && go build && cd ..
	cd messages/pg && rm -r go.sum || true && go mod tidy && go mod tidy && go build && cd ../..
	cd category && rm -r go.sum || true && go mod tidy && go build && cd ..
	cd category/pg && rm -r go.sum || true && go mod tidy && go build && cd ../..
	cd job && rm -r go.sum || true && go mod tidy && go build && cd ..
	cd job/pg && rm -r go.sum || true && go mod tidy && go build && cd ../..
	cd user && rm -r go.sum || true && go mod tidy && go build && cd ..
	cd user/pg && rm -r go.sum || true && go mod tidy && go build && cd ../..

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

