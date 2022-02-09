.PHONY:
.SILENT:

filename = "xookx"

build:
	# Linux
	GOOS=linux GOARCH=amd64 CGO_CFLAGS="-g -O2 -w" CGO_ENABLED=1 go build -ldflags "-s -H linux" -o ./.bin/$(filename)__linux-x86_64 cmd/main.go
	GOOS=linux GOARCH=386 CGO_CFLAGS="-g -O2 -w" CGO_ENABLED=1 go build -ldflags "-s -H linux" -o ./.bin/$(filename)__linux-x86 cmd/main.go

	# Linux ARM
#   GOOS=linux GOARCH=arm64 CGO_CFLAGS="-g -O2 -w" CGO_ENABLED=1 go build -ldflags "-s -H linux" -o ./.bin/$(filename)__linux-arm64 cmd/main.go
#	GOOS=linux GOARCH=arm CGO_CFLAGS="-g -O2 -w" CGO_ENABLED=1 go build -ldflags "-s -H linux" -o ./.bin/$(filename)__linux-arm32 cmd/main.go
#
#	# Windows
	GOOS=windows GOARCH=amd64 CGO_CFLAGS="-g -O2 -w" CGO_ENABLED=1 go build -ldflags "-s -H windows" -o ./.bin/$(filename)__windows-x86_64 cmd/main.go
	GOOS=windows GOARCH=386 CGO_CFLAGS="-g -O2 -w" CGO_ENABLED=1 go build -ldflags "-s -H windows" -o ./.bin/$(filename)__windows-x86 cmd/main.go
#
#	# Windows ARM
#	GOOS=windows GOARCH=arm64 CGO_CFLAGS="-g -O2 -w" CGO_ENABLED=1 go build -ldflags "-s -H windows" -o ./.bin/$(filename)__windows-arm64 cmd/main.go
#	GOOS=windows GOARCH=arm CGO_CFLAGS="-g -O2 -w" CGO_ENABLED=1 go build -ldflags "-s -H windows" -o ./.bin/$(filename)__windows-arm32 cmd/main.go

#	# Android x86
#	GOOS=android GOARCH=amd64 CGO_CFLAGS="-g -O2 -w" CGO_ENABLED=1 go build -ldflags "-s -H android" -o ./.bin/$(filename)__android-x86_64 cmd/main.go
#	GOOS=android GOARCH=386 CGO_CFLAGS="-g -O2 -w" CGO_ENABLED=1  go build -ldflags "-s -H android" -o ./.bin/$(filename)__android-x86 cmd/main.go

#   Android ARM
#	GOOS=android GOARCH=arm64 CGO_CFLAGS="-g -O2 -w" CGO_ENABLED=1 go build -ldflags "-s -H android" -o ./.bin/$(filename)__android-arm64 cmd/main.go
#	GOOS=android GOARCH=arm GOARM=7 CGO_CFLAGS="-g -O2 -w" CGO_ENABLED=1 go build -ldflags "-s -H android" -o ./.bin/$(filename)__android-arm32 cmd/main.go

	# MacOS x86
#	GOOS=darwin GOARCH=amd64 CGO_CFLAGS="-g -O2 -w" CGO_ENABLED=1 go build -ldflags "-s -H darwin" -o ./.bin/$(filename)__darwin-x86_64 cmd/main.go

	# MacOS ARM
#	GOOS=darwin GOARCH=arm64 CGO_CFLAGS="-g -O2 -w" CGO_ENABLED=1 go build -ldflags "-s -H darwin" -o ./.bin/$(filename)__darwin-arm64 cmd/main.go

run: build
	./.bin/$(filename)__linux-x86_64
