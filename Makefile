#
# Simple Makefile for conviently testing, building and deploying experiment.
#
export GOPATH=$(pwd)
export GOBIN="$HOME/bin"

test:
	go test

clean:
	if [ -f aspace ]; then rm aspace; fi

build: 
	go build
	go build cmds/aspace/aspace.go

install:
	if [ ! -d $GOBIN ]; then mkdir -p $GOBIN; fi
	go install cmds/aspace/aspace.go
