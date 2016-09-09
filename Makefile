install:
	go install -v

build:
	go build -v ./...

lint:
	golint ./...
	go vet ./...

test:
	go test -v ./... --cover

deps: dev-deps
	go get gopkg.in/redis.v3
	go get github.com/nats-io/nats
	go get github.com/ernestio/ernest-config-client
	go get github.com/tidwall/gjson

dev-deps:
	go get github.com/golang/lint/golint
	go get github.com/smartystreets/goconvey/convey

clean:
	go clean

dist-clean:
	rm -rf pkg src bin

