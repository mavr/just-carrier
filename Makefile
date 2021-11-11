PROJECTNAME = $(shell basename "$(PWD)")

BINDIRECTORY = bin

REVISIONBRANCH = $(shell git rev-parse --abbrev-ref HEAD)
REVISIONCOMMIT = $(shell git rev-parse HEAD | head -c 8)
REVISIONDATE = $(shell date +%Y.%m.%d-%H:%M:%S)
REVISIONVERSION = $(shell git describe --tags $(shell git rev-list --tags --max-count=1))
REVISION = $(REVISIONBRANCH)-$(REVISIONCOMMIT)-$(REVISIONDATE)

LDFLAGS = -ldflags="-X main.revision=$(REVISION) -X main.version=$(REVISIONVERSION)"

BUILDFLAGS = -v

build:
	CGO_ENABLED=0 go build $(LDFLAGS) $(BUILDFLAGS) -o $(BINDIRECTORY)/receiver ./cmd/receiver/*.go

modules:
	go mod download

vendor:
	go mod vendor

test:
	go test ./...

test.v:
	go test -v ./...

run: build
	$(BINDIRECTORY)/receiver

docker.run:
	docker-compose up -d receiver

docker.log:
	docker-compose logs -f

# docker.build:

# docker.deploy: docker.build

# mock.generate:
	