#
# Simple Makefile for conviently testing, building and deploying experiment.
#
PROG = cait

PROGRAM_liST = cait genpages sitemapper indexpages servepages 

build: api $(PROGRAM_LIST)

api: cait.go api.go export.go schema.go search.go views.go
	env CGO_ENABLED=0 go build

cait: api cmds/cait/cait.go
	env CGO_ENABLED=0 go build -o bin/cait cmds/cait/cait.go

genpages: api cmds/genpages/genpages.go
	env CGO_ENABLED=0 go build -o bin/genpages cmds/genpages/genpages.go

indexpages: api cmds/indexpages/indexpages.go
	env CGO_ENABLED=0 go build -o bin/indexpages cmds/indexpages/indexpages.go

servepages: api cmds/servepages/servepages.go
	env CGO_ENABLED=0 go build -o bin/servepages cmds/servepages/servepages.go

sitemapper: api cmds/sitemapper/sitemapper.go
	env CGO_ENABLED=0 go build -o bin/sitemapper cmds/sitemapper/sitemapper.go

test:
	go test

clean:
	if [ -d bin ]; then /bin/rm -fR bin; fi
	if [ -d dist ]; then /bin/rm -fR dist; fi
	if [ -f $(PROG)-release.zip ]; then /bin/rm $(PROG)-release.zip; fi

install:
	env CGO_ENABLED=0 GOBIN=$(HOME)/bin go install cmds/cait/cait.go
	env CGO_ENABLED=0 GOBIN=$(HOME)/bin go install cmds/genpages/genpages.go
	env CGO_ENABLED=0 GOBIN=$(HOME)/bin go install cmds/indexpages/indexpages.go
	env CGO_ENABLED=0 GOBIN=$(HOME)/bin go install cmds/servepages/servepages.go
	env CGO_ENABLED=0 GOBIN=$(HOME)/bin go install cmds/sitemapper/sitemapper.go

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

