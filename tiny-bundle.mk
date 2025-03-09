target=x86_64-linux-musl
export CC=zig cc -target $(target)
export CXX=zig c++ -target $(target)
export CGO_ENABLED=1
TAGS='static,osuergo,netgo'
EXTLDFLAGS="-static -Oz -s"
LDFLAGS='-linkmode=external -extldflags $(EXTLDFLAGS)'
build:
	@go build -tags $(TAGS) -ldflags $(LDFLAGS) -o bin/app_prod cmd/main.go
	@upx bin/app_prod
	@echo "compiled you application with all its assets to a single binary => bin/app_prod"