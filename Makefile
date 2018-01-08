.PHONY: clean test build

build:
	make -j3 gonelogin_amd64 gonelogin_darwin gonelogin.exe

clean:
	rm -f .get-deps
	rm -f *_amd64 *_darwin *.exe

test: .get-deps *.go
	go test -v *.go

.get-deps: *.go
	go get -t -d -v ./...
	touch .get-deps

gonelogin_amd64: .get-deps *.go
	GOOS=linux GOARCH=amd64 go build -o $@ *.go

gonelogin_darwin: .get-deps *.go
	GOOS=darwin go build -o $@ *.go

gonelogin.exe: .get-deps *.go
	GOOS=windows GOARCH=amd64 go build -o $@ *.go
