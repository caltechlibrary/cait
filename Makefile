#
# Simple Makefile for conviently testing, building and deploying experiment.
#
build: api.go  cait.go export.go  models.go  search.go  views.go
	go build
	go build -o bin/cait cmds/cait/cait.go
	go build -o bin/caitpage cmds/caitpage/caitpage.go
	go build -o bin/caitindexer cmds/caitindexer/caitindexer.go
	go build -o bin/caitserver cmds/caitserver/caitserver.go
	go build -o bin/xlsximporter cmds/xlsximporter/xlsximporter.go

test:
	go test

clean:
	if [ -f bin/cait ]; then rm bin/cait; fi
	if [ -f bin/caitpage ]; then rm bin/caitpage; fi
	if [ -f bin/caitindexer ]; then rm bin/caitindexer; fi
	if [ -f bin/caitserver ]; then rm bin/caitserver; fi
	if [ -f bin/xlsximporter ]; then rm bin/xlsximporter; fi
