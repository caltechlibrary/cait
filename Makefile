#
# Simple Makefile for conviently testing, building and deploying experiment.
#
export GOPATH=$(pwd)
export GOBIN=$HOME/bin

test:
	go test

clean:
	if [ -f aspace ]; then rm aspace; fi

build: 
	go build
	go build cmds/gospace/gospace.go

install:
	if [ ! -d $GOBIN ]; then mkdir -p $GOBIN; fi
	go install cmds/gospace/gospace.go
