#
# Simple Makefile for conviently testing, building and deploying experiment.
#
build: api.go  aspace.go export.go  models.go  search.go  views.go
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

