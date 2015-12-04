#!/bin/bash

function setupGolang {
    cd
    mkdir -p bin
    ORIGINAL_PATH="$PATH"
    mkdir -p src
    # Build a bootstrap version of Go
    git clone git@github.com:golang/go.git go1.4
    cd go1.4
    export GOBIN=$HOME/go1.4/bin
    git checkout go1.4.2
    cd src
    ./all.bash
    # Now build the current version of Go
    cd
    git clone git@github.com:golang/go.git go
    cd go/src
    export GOBIN=$HOME/go/bin
    ./all.bash
    # Update our local environment
    cd
    # Add a configuration examples to the .bashrc file
    echo 'You problably want to add the following to your .bashrc or .profile'
    echo ''
    echo '# Golang Setup '$(date)
    echo 'export PATH=$PATH:$HOME/bin:$HOME/go/bin'
    echo 'export GOPATH=$HOME'
    echo 'export GOBIN=$HOME/bin'
    echo
}

GO_CMD=$(which go)
if [ "$GO_CMD" = "" ]; then
    setupGolang
else
    echo "Go installed at $GO_CMD"
    echo "Version is "$(go version)
    echo "Gospace needs version 1.5 or better to compile"
fi
