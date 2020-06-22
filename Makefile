schema: schema.fbs
	flatc --go --go-namespace internal.flatbuffers schema.fbs

test:
	go test -v memrec/internal

build:
	go build -o build/memrec cmd/memrec.go

install:
	go install cmd/memrec.go

.PHONY: schema test build install