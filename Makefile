COMMIT=`git rev-parse --short HEAD`
LDFLAGS=-ldflags "-X main.gitVersion=${COMMIT}"

install:
	go install ${LDFLAGS} mod7
all: install
