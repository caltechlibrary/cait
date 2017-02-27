#
# Simple Makefile for conviently testing, building and deploying experiment.
#
PROJECT = cait

VERSION = $(shell grep -m 1 'Version =' $(PROJECT).go | cut -d\" -f 2)

BRANCH = $(shell git branch | grep '* ' | cut -d\  -f 2)

GOPATH = $(HOME)

PROGRAM_LIST = bin/cait bin/cait-genpages bin/cait-sitemapper bin/cait-indexpages bin/cait-servepages 

API = cait.go io.go export.go schema.go search.go views.go

CMDS = cmds/*/*.go

build: $(API) $(PROGRAM_LIST) $(CMDS)


api: $(API)
	env GOPATH=$(HOME) go build

cait: bin/cait

cait-genpages: bin/cait-genpages

cait-sitemapper: bin/cait-sitemapper

cait-indexpages: bin/cait-indexpages

cait-servepages: bin/cait-servepages

bin/cait: $(API) cmds/cait/cait.go
	env GOPATH=$(HOME) go build -o bin/cait cmds/cait/cait.go

bin/cait-genpages: $(API)  cmds/cait-genpages/cait-genpages.go
	env GOPATH=$(HOME) go build -o bin/cait-genpages cmds/cait-genpages/cait-genpages.go

bin/cait-indexpages: $(API) cmds/cait-indexpages/cait-indexpages.go
	env GOPATH=$(HOME) go build -o bin/cait-indexpages cmds/cait-indexpages/cait-indexpages.go

bin/cait-servepages: $(API) cmds/cait-servepages/cait-servepages.go
	env GOPATH=$(HOME) go build -o bin/cait-servepages cmds/cait-servepages/cait-servepages.go

bin/cait-sitemapper: $(API) cmds/cait-sitemapper/cait-sitemapper.go
	env GOPATH=$(HOME) go build -o bin/cait-sitemapper cmds/cait-sitemapper/cait-sitemapper.go

test:
	env GOPATH=$(HOME) go test

clean:
	if [ -d bin ]; then /bin/rm -fR bin; fi
	if [ -d dist ]; then /bin/rm -fR dist; fi
	if [ -f $(PROJECT)-$(VERSION)-release.zip ]; then /bin/rm $(PROJECT)-$(VERSION)-release.zip; fi

install:
	env GOPATH=$(HOME) GOBIN=$(HOME)/bin go install cmds/cait/cait.go
	env GOPATH=$(HOME) GOBIN=$(HOME)/bin go install cmds/cait-genpages/cait-genpages.go
	env GOPATH=$(HOME) GOBIN=$(HOME)/bin go install cmds/cait-indexpages/cait-indexpages.go
	env GOPATH=$(HOME) GOBIN=$(HOME)/bin go install cmds/cait-servepages/cait-servepages.go
	env GOPATH=$(HOME) GOBIN=$(HOME)/bin go install cmds/cait-sitemapper/cait-sitemapper.go

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
	env GOOS=linux GOARCH=amd64 go build -o dist/linux-amd64/cait cmds/cait/cait.go
	env GOOS=linux GOARCH=amd64 go build -o dist/linux-amd64/cait-genpages cmds/cait-genpages/cait-genpages.go
	env GOOS=linux GOARCH=amd64 go build -o dist/linux-amd64/cait-sitemapper cmds/cait-sitemapper/cait-sitemapper.go
	env GOOS=linux GOARCH=amd64 go build -o dist/linux-amd64/cait-indexpages cmds/cait-indexpages/cait-indexpages.go
	env GOOS=linux GOARCH=amd64 go build -o dist/linux-amd64/cait-servepages cmds/cait-servepages/cait-servepages.go

dist/windows-amd64: *.go cmds/cait/cait.go cmds/cait-genpages/cait-genpages.go cmds/cait-sitemapper/cait-sitemapper.go cmds/cait-indexpages/cait-indexpages.go cmds/cait-servepages/cait-servepages.go
	env GOOS=windows GOARCH=amd64 go build -o dist/windows-amd64/cait cmds/cait/cait.go
	env GOOS=windows GOARCH=amd64 go build -o dist/windows-amd64/cait-genpages cmds/cait-genpages/cait-genpages.go
	env GOOS=windows GOARCH=amd64 go build -o dist/windows-amd64/cait-sitemapper cmds/cait-sitemapper/cait-sitemapper.go
	env GOOS=windows GOARCH=amd64 go build -o dist/windows-amd64/cait-indexpages cmds/cait-indexpages/cait-indexpages.go
	env GOOS=windows GOARCH=amd64 go build -o dist/windows-amd64/cait-servepages cmds/cait-servepages/cait-servepages.go

dist/macosx-amd64: *.go cmds/cait/cait.go cmds/cait-genpages/cait-genpages.go cmds/cait-sitemapper/cait-sitemapper.go cmds/cait-indexpages/cait-indexpages.go cmds/cait-servepages/cait-servepages.go
	env GOOS=darwin GOARCH=amd64 go build -o dist/macosx-amd64/cait cmds/cait/cait.go
	env GOOS=darwin GOARCH=amd64 go build -o dist/macosx-amd64/cait-genpages cmds/cait-genpages/cait-genpages.go
	env GOOS=darwin GOARCH=amd64 go build -o dist/macosx-amd64/cait-sitemapper cmds/cait-sitemapper/cait-sitemapper.go
	env GOOS=darwin GOARCH=amd64 go build -o dist/macosx-amd64/cait-indexpages cmds/cait-indexpages/cait-indexpages.go
	env GOOS=darwin GOARCH=amd64 go build -o dist/macosx-amd64/cait-servepages cmds/cait-servepages/cait-servepages.go

dist/raspbian-arm7: *.go cmds/cait/cait.go cmds/cait-genpages/cait-genpages.go cmds/cait-sitemapper/cait-sitemapper.go cmds/cait-indexpages/cait-indexpages.go cmds/cait-servepages/cait-servepages.go
	env GOOS=linux GOARCH=arm GOARM=7 go build -o dist/raspbian-arm7/cait cmds/cait/cait.go
	env GOOS=linux GOARCH=arm GOARM=7 go build -o dist/raspbian-arm7/cait-genpages cmds/cait-genpages/cait-genpages.go
	env GOOS=linux GOARCH=arm GOARM=7 go build -o dist/raspbian-arm7/cait-sitemapper cmds/cait-sitemapper/cait-sitemapper.go
	env GOOS=linux GOARCH=arm GOARM=7 go build -o dist/raspbian-arm7/cait-indexpages cmds/cait-indexpages/cait-indexpages.go
	env GOOS=linux GOARCH=arm GOARM=7 go build -o dist/raspbian-arm7/cait-servepages cmds/cait-servepages/cait-servepages.go

dist/raspbian-arm6: *.go cmds/cait/cait.go cmds/cait-genpages/cait-genpages.go cmds/cait-sitemapper/cait-sitemapper.go cmds/cait-indexpages/cait-indexpages.go cmds/cait-servepages/cait-servepages.go
	env GOOS=linux GOARCH=arm GOARM=6 go build -o dist/raspbian-arm6/cait cmds/cait/cait.go
	env GOOS=linux GOARCH=arm GOARM=6 go build -o dist/raspbian-arm6/cait-genpages cmds/cait-genpages/cait-genpages.go
	env GOOS=linux GOARCH=arm GOARM=6 go build -o dist/raspbian-arm6/cait-sitemapper cmds/cait-sitemapper/cait-sitemapper.go
	env GOOS=linux GOARCH=arm GOARM=6 go build -o dist/raspbian-arm6/cait-indexpages cmds/cait-indexpages/cait-indexpages.go
	env GOOS=linux GOARCH=arm GOARM=6 go build -o dist/raspbian-arm6/cait-servepages cmds/cait-servepages/cait-servepages.go

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


