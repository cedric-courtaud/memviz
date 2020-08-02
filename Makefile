schema: schema.fbs
	flatc --go --go-namespace internal.flatbuffers schema.fbs

test:
	go test -v github.com/cedric-courtaud/memviz/internal

build:
	go build -o build/memviz cmd/memrec.go

install:
	go install cmd/memrec.go

.PHONY: schema test build install