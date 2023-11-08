BINARY_NAME=cryptid-mqtt-broker

all: deps build

deps:
	go get -d ./...

build:
	go build -o ${BINARY_NAME} main.go

install:
	install -v ${BINARY_NAME} /usr/bin/${BINARY_NAME}
	install -v systemd.service /etc/systemd/system/${BINARY_NAME}.service

run:
	go build -o ${BINARY_NAME} main.go
	./${BINARY_NAME}

test:
	go test -v main.go

clean:
	go clean
	rm ${BINARY_NAME}
