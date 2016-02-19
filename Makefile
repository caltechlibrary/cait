#
# Simple Makefile for conviently testing, building and deploying experiment.
#
build: api.go cait.go export.go models.go search.go views.go js.go
	go build
	go build -o bin/cait cmds/cait/cait.go
	go build -o bin/caitjs cmds/caitjs/caitjs.go
	go build -o bin/caitserver cmds/caitserver/caitserver.go
	go build -o bin/genpages cmds/genpages/genpages.go
	go build -o bin/indexpages cmds/indexpages/indexpages.go
	go build -o bin/xlsximporter cmds/xlsximporter/xlsximporter.go
	go build -o bin/sitemapper cmds/sitemapper/sitemapper.go
	go build -o bin/indexdataset cmds/indexdataset/indexdataset.go
	go build -o bin/searchdataset cmds/searchdataset/searchdataset.go

test:
	go test

clean:
	if [ -d bin ]; then rm bin/*; fi
