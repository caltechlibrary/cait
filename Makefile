#
# Simple Makefile for conviently testing, building and deploying experiment.
#
build: aspace.go models.go views.go export.go
	go build
	go build -o bin/aspace cmds/aspace/aspace.go
	go build -o bin/aspacepage cmds/aspacepage/aspacepage.go
	go build -o bin/aspaceindexer cmds/aspaceindexer/aspaceindexer.go
	go build -o bin/aspacesearch cmds/aspacesearch/aspacesearch.go
	go build -o bin/xlsximporter cmds/xlsximporter/xlsximporter.go

test:
	go test

clean:
	if [ -f bin/aspace ]; then rm bin/aspace; fi
	if [ -f bin/aspacepage ]; then rm bin/aspacepage; fi
	if [ -f bin/aspaceindexer ]; then rm bin/aspaceindexer; fi
	if [ -f bin/aspacesearch ]; then rm bin/aspacesearch; fi
	if [ -f bin/xlsximporter ]; then rm bin/xlsximporter; fi

install:
	if [ ! -d "$GOBIN" ] && [ "$GOBIN" != "" ]; then mkdir -p "$GOBIN"; fi
	go install cmds/aspace/aspace.go
	go install cmds/aspacepage/aspacepage.go
	go install cmds/aspaceindexer/aspaceindexer.go
	go install cmds/aspacesearch/aspacesearch.go
	go install cmds/xlsximporter/xlsximporter.go
