#
# Simple Makefile for conviently testing, building and deploying experiment.
#
build: api.go cait.go export.go schema.go search.go views.go js.go
	go build
	go build -o bin/cait cmds/cait/cait.go
	go build -o bin/genpages cmds/genpages/genpages.go
	go build -o bin/indexpages cmds/indexpages/indexpages.go
	go build -o bin/servepages cmds/servepages/servepages.go
	go build -o bin/sitemapper cmds/sitemapper/sitemapper.go
	go build -o bin/xlsximporter cmds/xlsximporter/xlsximporter.go

test:
	go test

clean:
	if [ -d bin ]; then rm bin/*; fi
