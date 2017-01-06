
# Installation

This is generalized instructions for a release.  For deployment suggestions see NOTES.md

## Compiled version

*cait* is a set of command line programs run from a shell like Bash. If you download the repository a compiled version is in the dist directory. The compiled binary matching your computer type and operating system can be copied to a bin directory in your PATH.

Compiled versions are available for Mac OS X (amd64 processor), Linux (amd64), Windows (amd64) and Rapsberry Pi (both ARM6 and ARM7)

VERSION_NUMBER is a [symantic version number](http://semver.org/) (e.g. v0.1.2)

### Mac OS X

1. Download **cait-VERSION_NUMBER-release.zip** from [https://github.com/caltechlibrary/cait/releases/latest](https://github.com/caltechlibrary/cait/releases/latest)
2. Open a finder window, find and unzip **cait-VERSION_NUMBER-release.zip**
3. Look in the unziped folder and find the files in *dist/macosx-amd64/*
4. Drag (or copy) both the programs (e.g. *cait*, *genpages*) to a "bin" directory in your path
5. Open and "Terminal" and run `cait -h` to confirm you were successful

### Windows

1. Download **cait-VERSION_NUMBER-release.zip** from [https://github.com/caltechlibrary/cait/releases/latest](https://github.com/caltechlibrary/cait/releases/latest)
2. Open the file manager find and unzip **cait-VERSION_NUMBER-release.zip**
3. Look in the unziped folder and find the files in *dist/windows-amd64/*
4. Drag (or copy) the programs (e.g. *cait.exe*, *genpages.exe*) to a "bin" directory in your path
5. Open Bash and and run `cait -h` to confirm you were successful

### Linux

1. Download **cait-VERSION_NUMBER-release.zip** from [https://github.com/caltechlibrary/cait/releases/latest](https://github.com/caltechlibrary/cait/releases/latest)
2. Find and unzip **cait-VERSION_NUMBER-release.zip**
3. In the unziped directory and find the files in *dist/linux-amd64/*
4. Copy the programs (e.g. *cait*, *genpages*) to a "bin" directory (e.g. cp ~/Downloads/cait-VERSION_NUMBER-release/dist/linux-amd64/cait ~/bin/)
5. From the shell prompt run `cait -h` to confirm you were successful

### Raspberry Pi

If you are using a Raspberry Pi 2 or later use the ARM7 VERSION_NUMBER, ARM6 is only for the first generaiton Raspberry Pi.

1. Download **cait-VERSION_NUMBER-release.zip** from [https://github.com/caltechlibrary/cait/releases/latest](https://github.com/caltechlibrary/cait/releases/latest)
2. Find and unzip **cait-VERSION_NUMBER-release.zip**
3. In the unziped directory and find the files in *dist/rasbian-arm7/*
4. Copy the programs to a "bin" directory (e.g. cp ~/Downloads/cait-VERSION_NUMBER-release/dist/rasbian-arm7/cait ~/bin/)
    + if you are using an original Raspberry Pi you should copy the ARM6 version instead (i.e. dist/raspbian-arm6)
5. From the shell prompt run `cait -h` to confirm you were successful


## Compiling from source

```shell
    go get github.com/blevesearch/bleve/...
    cd src/github.com/blevesearch/belve
    git checkout v0.5.0
    cd
    go get github.com/caltechlibrary/cli
    go get github.com/caltechlibrary/tmplfn
    go get github.com/caltechlibrary/cait
    cd src/github.com/caltechlibrary/cait
    make
    make test
    make install
```

