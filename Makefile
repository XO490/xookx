.PHONY: all build clean
.SILENT:

filename = "xookx"
binpath = "bin"
maingo = "cmd/main.go"

mkdir:
	mkdir -p $(binpath)

build: mkdir
#	 Linux
	GOOS=linux GOARCH=amd64 CGO_ENABLED=1 go build -gcflags "-m" -ldflags "-s" -o $(binpath)/$(filename)-l64 $(maingo)

#	Windows
	GOOS=windows GOARCH=amd64 CC="x86_64-w64-mingw32-gcc" CGO_ENABLED=1 go build -ldflags "-s" -o $(binpath)/$(filename)-w64.exe $(maingo)

run: build
	./bin/$(filename)-l64
