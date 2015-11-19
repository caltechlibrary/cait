#!/bin/bash

function setupGolang {
    cd
    mkdir -p bin
    ORIGINAL_PATH="$PATH"
    mkdir -p src
    # Build a bootstrap version of Go
    git clone https://github.com/golang/go.git go1.4
    cd go1.4
    export GOBIN=$HOME/go1.4/bin
    git checkout go1.4.2
    cd src
    ./all.bash
    # Now build the current version of Go
    cd
    git clone https://github.com/golang/go.git go
    cd go/src
    export GOBIN=$HOME/go/bin
    ./all.bash
    # Update our local environment
    cd
    # Add a configuration examples to the .bashrc file
    echo '## Golang Setup '$(date) >> .bashrc
    echo '#export PATH=$PATH:$HOME/bin:$HOME/go/bin' >> .bashrc
    echo '#export GOPATH=$HOME' >> .bashrc
    echo '#export GOBIN=$HOME/bin' >> .bashrc
}

function setupGoSpace {
    cd
    mkdir -p src
    mkdir -p $HOME/bin
    cd src
    git clone https://github.com/caltechlibrary/gospace.git
    cd gospace
    export GOPATH=$HOME
    export GOBIN=$HOME/bin
    go test && go install
}

GO_CMD=$(which go)
if [ "$GO_CMD" = "" ]; then
    setupGolang
fi
GOSPACE_CMD=$(which aspace)
if [ "$GOSPACE_CMD" = "" ]; then
    setupGoSpace
fi
