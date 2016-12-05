#!/bin/bash
#
# Make releases for Linux/amd64, Linux/ARM7 (Raspberry Pi), Windows, and Mac OX X (darwin)
#
RELEASE_NAME=cait

PROGRAM_LIST="cait genpages sitemapper indexpages servepages"

for PROGNAME in $PROGRAM_LIST; do
    echo "Compiling $PROGNAME for all architectures"
    env CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o dist/linux-amd64/$PROGNAME cmds/$PROGNAME/$PROGNAME.go
    env CGO_ENABLED=0 GOOS=darwin	GOARCH=amd64 go build -o dist/macosx-amd64/$PROGNAME cmds/$PROGNAME/$PROGNAME.go
    env CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -o dist/windows-amd64/$PROGNAME.exe cmds/$PROGNAME/$PROGNAME.go
    env CGO_ENABLED=0 GOOS=linux GOARCH=arm GOARM=6 go build -o dist/raspberrypi-arm6/$PROGNAME cmds/$PROGNAME/$PROGNAME.go
    env CGO_ENABLED=0 GOOS=linux GOARCH=arm GOARM=7 go build -o dist/raspberrypi-arm7/$PROGNAME cmds/$PROGNAME/$PROGNAME.go
done
echo "Zipping release $RELEASE_NAME"
zip -r $RELEASE_NAME-release.zip README.md INSTALL.md NOTES.md LICENSE etc/cait.bash-example etc/systemd/system/*example scripts/harvest-*.bash templates/* dist/*

