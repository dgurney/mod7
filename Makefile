COMMIT=`git rev-parse --short HEAD`
BRANCH=`git rev-parse --abbrev-ref HEAD`
LDFLAGS=-ldflags "-w -X main.gitCommit=${COMMIT} -X main.gitBranch=${BRANCH}"
PROGRAM=mod7

install:
	go install ${LDFLAGS} ${PROGRAM}
windows:
	GOOS=windows GOARCH=amd64 go build ${LDFLAGS} -o build/windows/amd64/${PROGRAM}.exe ${PROGRAM} 
	GOOS=windows GOARCH=386 go build ${LDFLAGS} -o build/windows/386/${PROGRAM}.exe ${PROGRAM} 
darwin:
	GOOS=darwin GOARCH=amd64 go build ${LDFLAGS} -o build/darwin/amd64/${PROGRAM} ${PROGRAM} 
	GOOS=darwin GOARCH=386 go build ${LDFLAGS} -o build/darwin/386/${PROGRAM} ${PROGRAM} 
freebsd:
	GOOS=freebsd GOARCH=amd64 go build ${LDFLAGS} -o build/freebsd/amd64/${PROGRAM} ${PROGRAM} 
	GOOS=freebsd GOARCH=386 go build ${LDFLAGS} -o build/freebsd/386/${PROGRAM} ${PROGRAM} 
linux:
	GOOS=linux GOARCH=amd64 go build ${LDFLAGS} -o build/linux/amd64/${PROGRAM} ${PROGRAM} 
	GOOS=linux GOARCH=arm64 go build ${LDFLAGS} -o build/linux/arm64/${PROGRAM} ${PROGRAM} 
	GOOS=linux GOARCH=arm GOARM=7 go build ${LDFLAGS} -o build/linux/armv7/${PROGRAM} ${PROGRAM} 
	GOOS=linux GOARCH=386 go build ${LDFLAGS} -o build/linux/386/${PROGRAM} ${PROGRAM} 
clean:
	rm -rf build/
test:
	GOCACHE=off	go test mod7/oem -cover
	GOCACHE=off	go test mod7/tendigit -cover
	GOCACHE=off	go test mod7/validation -cover
bench:
	go test -run=sonic mod7/oem -bench=.
	go test -run=sonic mod7/tendigit -bench=.
	go test -run=sonic mod7/validation -bench=.
cross: windows darwin freebsd linux
all: install cross
