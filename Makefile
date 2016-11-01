#
# Simple Makefile for conviently testing, building and deploying experiment.
#
PROG = cait

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
	if [ -d bin ]; then /bin/rm -fR bin; fi
	if [ -d dist ]; then /bin/rm -fR dist; fi
	if [ -f $(PROG)-binary-release.zip ]; then /bin/rm $(PROG)-binary-release.zip; fi

install:
	env GOBIN=$(HOME)/bin go install cmds/cait/cait.go
	env GOBIN=$(HOME)/bin go install cmds/genpages/genpages.go
	env GOBIN=$(HOME)/bin go install cmds/indexpages/indexpages.go
	env GOBIN=$(HOME)/bin go install cmds/servepages/servepages.go
	env GOBIN=$(HOME)/bin go install cmds/sitemapper/sitemapper.go
	env GOBIN=$(HOME)/bin go install cmds/xlsximporter/xlsximporter.go

website:
	./mk-website.bash

save:
	./mk-website.bash
	git commit -am "Quick save"
	git push origin master

publish:
	./mk-website.bash
	./publish.bash

release:
	./mk-release.bash

