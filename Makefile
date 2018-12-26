COMMIT=`git rev-parse --short HEAD`
LDFLAGS=-ldflags "-w -X main.gitVersion=${COMMIT}"
PROGRAM=mod7

install:
	go install ${LDFLAGS} mod7
windows:
	GOOS=windows GOARCH=amd64 go build ${LDFLAGS} -o build/windows/amd64/${PROGRAM}.exe mod7 
	GOOS=windows GOARCH=386 go build ${LDFLAGS} -o build/windows/386/${PROGRAM}.exe mod7 
darwin:
	GOOS=darwin GOARCH=amd64 go build ${LDFLAGS} -o build/darwin/amd64/${PROGRAM} mod7 
	GOOS=darwin GOARCH=386 go build ${LDFLAGS} -o build/darwin/386/${PROGRAM} mod7 
freebsd:
	GOOS=freebsd GOARCH=amd64 go build ${LDFLAGS} -o build/freebsd/amd64/${PROGRAM} mod7 
	GOOS=freebsd GOARCH=386 go build ${LDFLAGS} -o build/freebsd/386/${PROGRAM} mod7 
linux:
	GOOS=linux GOARCH=amd64 go build ${LDFLAGS} -o build/linux/amd64/${PROGRAM} mod7 
	GOOS=linux GOARCH=arm64 go build ${LDFLAGS} -o build/linux/arm64/${PROGRAM} mod7 
	GOOS=linux GOARCH=386 go build ${LDFLAGS} -o build/linux/386/${PROGRAM} mod7 
clean:
	rm -rf build/
cross: windows darwin freebsd linux
all: install cross
