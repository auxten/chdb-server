.PHONY: update_libchdb all test clean

update_libchdb:
	./update_libchdb.sh

install:
	curl -sL https://lib.chdb.io | bash

test:
	CGO_ENABLED=1 go test -v -coverprofile=coverage.out ./...

run:
	# add ld path for linux and macos
	LD_LIBRARY_PATH=/usr/local/lib DYLD_LIBRARY_PATH=/usr/local/lib CGO_ENABLED=1 go run -ldflags '-extldflags "-Wl,-rpath,/usr/local/lib"' main.go

build:
	CGO_ENABLED=1 go build -ldflags '-extldflags "-Wl,-rpath,/usr/local/lib"' -o chdb-server main.go
	install_name_tool -change libchdb.so /usr/local/lib/libchdb.so chdb-server
