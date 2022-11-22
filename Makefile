build:
	go build -ldflags "-X main.versionSHA=`git rev-parse --short HEAD`"