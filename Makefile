schema: schema.fbs
	flatc --go --go-namespace internal.flatbuffers schema.fbs

test:
	go test -v memrec/internal



.PHONY: schema