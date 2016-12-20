#
# Simple Makefile for conviently testing, building and deploying experiment.
#
PROG = cait

VERSION = $(shell grep -m 1 'Version =' cait.go | cut -d\" -f 2)

BRANCH = $(shell git branch | grep '* ' | cut -d\  -f 2)

PROGRAM_LIST = bin/cait bin/genpages bin/sitemapper bin/indexpages bin/servepages 

API = cait.go api.go export.go schema.go search.go views.go

CMDS = cmds/*/*.go

build: $(API) $(PROGRAM_LIST) $(CMDS)

api: $(API)
	env CGO_ENABLED=0 go build

cait: bin/cait

genpages: bin/genpages

sitemapper: bin/sitemapper

indexpages: bin/indexpages

servepages: bin/servepages

bin/cait: $(API) cmds/cait/cait.go
	env CGO_ENABLED=0 go build -o bin/cait cmds/cait/cait.go

bin/genpages: $(API)  cmds/genpages/genpages.go
	env CGO_ENABLED=0 go build -o bin/genpages cmds/genpages/genpages.go

bin/indexpages: $(API) cmds/indexpages/indexpages.go
	env CGO_ENABLED=0 go build -o bin/indexpages cmds/indexpages/indexpages.go

bin/servepages: $(API) cmds/servepages/servepages.go
	env CGO_ENABLED=0 go build -o bin/servepages cmds/servepages/servepages.go

bin/sitemapper: $(API) cmds/sitemapper/sitemapper.go
	env CGO_ENABLED=0 go build -o bin/sitemapper cmds/sitemapper/sitemapper.go

test:
	go test

clean:
	if [ -d bin ]; then /bin/rm -fR bin; fi
	if [ -d dist ]; then /bin/rm -fR dist; fi
	if [ -f $(PROG)-$(VERSION)-release.zip ]; then /bin/rm $(PROG)-$(VERSION)-release.zip; fi

install:
	env CGO_ENABLED=0 GOBIN=$(HOME)/bin go install cmds/cait/cait.go
	env CGO_ENABLED=0 GOBIN=$(HOME)/bin go install cmds/genpages/genpages.go
	env CGO_ENABLED=0 GOBIN=$(HOME)/bin go install cmds/indexpages/indexpages.go
	env CGO_ENABLED=0 GOBIN=$(HOME)/bin go install cmds/servepages/servepages.go
	env CGO_ENABLED=0 GOBIN=$(HOME)/bin go install cmds/sitemapper/sitemapper.go

website:
	./mk-website.bash

save:
	git commit -am "Quick save"
	git push origin $(BRANCH)

publish:
	./mk-website.bash
	./publish.bash

release:
	./mk-release.bash

