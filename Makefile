GOFILES := $(shell find . -type f -name '*.go')

run: calsync
	./calsync

calsync: $(GOFILES)
	go build -o calsync ./main

.PHONY: run
