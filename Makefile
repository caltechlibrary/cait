#
# Simple Makefile for conviently testing, building and deploying experiment.
#
PROJECT = cait

VERSION = $(shell grep -m 1 'Version =' $(PROJECT).go | cut -d\" -f 2)

BRANCH = $(shell git branch | grep '* ' | cut -d\  -f 2)

PROGRAM_LIST = bin/cait bin/cait-genpages bin/cait-sitemapper bin/cait-indexpages bin/cait-servepages 

API = cait.go api.go io.go export.go schema.go search.go views.go

CMDS = cmds/*/*.go

build: dependency $(API) $(PROGRAM_LIST) $(CMDS)

dependency: cli tmplfn bleve

cli: $(GOPATH)/src/github.com/caltechlibrary/cli/cli.go 

$(GOPATH)/src/github.com/caltechlibrary/cli/cli.go:
	go get github.com/caltechlibrary/cli

tmplfn: $(GOPATH)/src/github.com/caltechlibrary/tmplfn/tmplfn.go

$(GOPATH)/src/github.com/caltechlibrary/tmplfn/tmplfn.go:
	go get github.com/caltechlibrary/tmplfn

bleve: $(GOPATH)/src/github.com/blevesearch/bleve/index.go

$(GOPATH)/src/github.com/blevesearch/bleve/index.go:
	go get github.com/blevesearch/bleve/...
	cd $(GOPATH)/src/github.com/blevesearch/bleve && git checkout v0.5.0

api: $(API)
	env CGO_ENABLED=0 go build

cait: bin/cait

cait-genpages: bin/cait-genpages

cait-sitemapper: bin/cait-sitemapper

cait-indexpages: bin/cait-indexpages

cait-servepages: bin/cait-servepages

bin/cait: $(API) cmds/cait/cait.go
	env CGO_ENABLED=0 go build -o bin/cait cmds/cait/cait.go

bin/cait-genpages: $(API)  cmds/cait-genpages/cait-genpages.go
	env CGO_ENABLED=0 go build -o bin/cait-genpages cmds/cait-genpages/cait-genpages.go

bin/cait-indexpages: $(API) cmds/cait-indexpages/cait-indexpages.go
	env CGO_ENABLED=0 go build -o bin/cait-indexpages cmds/cait-indexpages/cait-indexpages.go

bin/cait-servepages: $(API) cmds/cait-servepages/cait-servepages.go
	env CGO_ENABLED=0 go build -o bin/cait-servepages cmds/cait-servepages/cait-servepages.go

bin/cait-sitemapper: $(API) cmds/cait-sitemapper/cait-sitemapper.go
	env CGO_ENABLED=0 go build -o bin/cait-sitemapper cmds/cait-sitemapper/cait-sitemapper.go

test:
	go test

clean:
	if [ -d bin ]; then /bin/rm -fR bin; fi
	if [ -d dist ]; then /bin/rm -fR dist; fi
	if [ -f $(PROJECT)-$(VERSION)-release.zip ]; then /bin/rm $(PROJECT)-$(VERSION)-release.zip; fi

install:
	env CGO_ENABLED=0 GOBIN=$(HOME)/bin go install cmds/cait/cait.go
	env CGO_ENABLED=0 GOBIN=$(HOME)/bin go install cmds/cait-genpages/cait-genpages.go
	env CGO_ENABLED=0 GOBIN=$(HOME)/bin go install cmds/cait-indexpages/cait-indexpages.go
	env CGO_ENABLED=0 GOBIN=$(HOME)/bin go install cmds/cait-servepages/cait-servepages.go
	env CGO_ENABLED=0 GOBIN=$(HOME)/bin go install cmds/cait-sitemapper/cait-sitemapper.go

website:
	./mk-website.bash

save:
	git commit -am "Quick save"
	git push origin $(BRANCH)

refresh:
	git fetch origin
	git pull origin $(BRANCH)

status:
	git status

publish:
	./mk-website.bash
	./publish.bash

dist/linux-amd64: *.go cmds/cait/cait.go cmds/cait-genpages/cait-genpages.go cmds/cait-sitemapper/cait-sitemapper.go cmds/cait-indexpages/cait-indexpages.go cmds/cait-servepages/cait-servepages.go
	env CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o dist/linux-amd64/cait cmds/cait/cait.go
	env CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o dist/linux-amd64/cait-genpages cmds/cait-genpages/cait-genpages.go
	env CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o dist/linux-amd64/cait-sitemapper cmds/cait-sitemapper/cait-sitemapper.go
	env CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o dist/linux-amd64/cait-indexpages cmds/cait-indexpages/cait-indexpages.go
	env CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o dist/linux-amd64/cait-servepages cmds/cait-servepages/cait-servepages.go

dist/windows-amd64: *.go cmds/cait/cait.go cmds/cait-genpages/cait-genpages.go cmds/cait-sitemapper/cait-sitemapper.go cmds/cait-indexpages/cait-indexpages.go cmds/cait-servepages/cait-servepages.go
	env CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -o dist/windows-amd64/cait cmds/cait/cait.go
	env CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -o dist/windows-amd64/cait-genpages cmds/cait-genpages/cait-genpages.go
	env CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -o dist/windows-amd64/cait-sitemapper cmds/cait-sitemapper/cait-sitemapper.go
	env CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -o dist/windows-amd64/cait-indexpages cmds/cait-indexpages/cait-indexpages.go
	env CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -o dist/windows-amd64/cait-servepages cmds/cait-servepages/cait-servepages.go

dist/macosx-amd64: *.go cmds/cait/cait.go cmds/cait-genpages/cait-genpages.go cmds/cait-sitemapper/cait-sitemapper.go cmds/cait-indexpages/cait-indexpages.go cmds/cait-servepages/cait-servepages.go
	env CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -o dist/macosx-amd64/cait cmds/cait/cait.go
	env CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -o dist/macosx-amd64/cait-genpages cmds/cait-genpages/cait-genpages.go
	env CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -o dist/macosx-amd64/cait-sitemapper cmds/cait-sitemapper/cait-sitemapper.go
	env CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -o dist/macosx-amd64/cait-indexpages cmds/cait-indexpages/cait-indexpages.go
	env CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -o dist/macosx-amd64/cait-servepages cmds/cait-servepages/cait-servepages.go

dist/raspbian-arm7: *.go cmds/cait/cait.go cmds/cait-genpages/cait-genpages.go cmds/cait-sitemapper/cait-sitemapper.go cmds/cait-indexpages/cait-indexpages.go cmds/cait-servepages/cait-servepages.go
	env CGO_ENABLED=0 GOOS=linux GOARCH=arm GOARM=7 go build -o dist/raspbian-arm7/cait cmds/cait/cait.go
	env CGO_ENABLED=0 GOOS=linux GOARCH=arm GOARM=7 go build -o dist/raspbian-arm7/cait-genpages cmds/cait-genpages/cait-genpages.go
	env CGO_ENABLED=0 GOOS=linux GOARCH=arm GOARM=7 go build -o dist/raspbian-arm7/cait-sitemapper cmds/cait-sitemapper/cait-sitemapper.go
	env CGO_ENABLED=0 GOOS=linux GOARCH=arm GOARM=7 go build -o dist/raspbian-arm7/cait-indexpages cmds/cait-indexpages/cait-indexpages.go
	env CGO_ENABLED=0 GOOS=linux GOARCH=arm GOARM=7 go build -o dist/raspbian-arm7/cait-servepages cmds/cait-servepages/cait-servepages.go

dist/raspbian-arm6: *.go cmds/cait/cait.go cmds/cait-genpages/cait-genpages.go cmds/cait-sitemapper/cait-sitemapper.go cmds/cait-indexpages/cait-indexpages.go cmds/cait-servepages/cait-servepages.go
	env CGO_ENABLED=0 GOOS=linux GOARCH=arm GOARM=6 go build -o dist/raspbian-arm6/cait cmds/cait/cait.go
	env CGO_ENABLED=0 GOOS=linux GOARCH=arm GOARM=6 go build -o dist/raspbian-arm6/cait-genpages cmds/cait-genpages/cait-genpages.go
	env CGO_ENABLED=0 GOOS=linux GOARCH=arm GOARM=6 go build -o dist/raspbian-arm6/cait-sitemapper cmds/cait-sitemapper/cait-sitemapper.go
	env CGO_ENABLED=0 GOOS=linux GOARCH=arm GOARM=6 go build -o dist/raspbian-arm6/cait-indexpages cmds/cait-indexpages/cait-indexpages.go
	env CGO_ENABLED=0 GOOS=linux GOARCH=arm GOARM=6 go build -o dist/raspbian-arm6/cait-servepages cmds/cait-servepages/cait-servepages.go

release: dist/linux-amd64 dist/windows-amd64 dist/macosx-amd64 dist/raspbian-arm7 dist/raspbian-arm6
	mkdir -p dist
	mkdir -p dist/etc/systemd/system
	mkdir -p dist/scripts
	cp -v README.md dist/
	cp -v LICENSE dist/
	cp -v INSTALL.md dist/
	cp -v NOTES.md dist/
	cp -vR templates dist/
	cp -vR scripts/harvest-*.bash dist/scripts/
	cp -vR etc/*-example dist/etc/
	cp -vR etc/systemd/system/*-example dist/etc/systemd/system/
	zip -r $(PROJECT)-$(VERSION)-release.zip dist/*


