#
# Simple Makefile for conviently testing, building and deploying experiment.
#
build: aspace.go models.go
	go build
	go build -o bin/aspace cmds/aspace/aspace.go
	go build -o bin/aspaceindexer cmds/aspaceindexer/aspaceindexer.go
	go build -o bin/aspacesearch cmds/aspacesearch/aspacesearch.go
	go build -o bin/aspacedashboard cmds/aspacedashboard/aspacedashboard.go

test:
	go test

clean:
	if [ -f bin/aspace ]; then rm bin/aspace; fi
	if [ -f bin/aspaceindexer ]; then rm bin/aspaceindexer; fi
	if [ -f bin/aspacesearch ]; then rm bin/aspacesearch; fi
	if [ -f bin/aspacedashboard ]; then rm bin/aspacedashboard; fi

install:
	if [ ! -d $GOBIN ] && [ "$GOBIN" != "" ]; then mkdir -p $GOBIN; fi
	go install cmds/aspace/aspace.go
	go install cmds/aspaceindexer/aspaceindexer.go
	go install cmds/aspacesearch/aspacesearch.go
	go install cmds/aspacedashboard/aspacedashboard.go
