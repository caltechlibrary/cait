#!/bin/bash

#
# Example webhook for updating archives.example.edu
#
cd /Sites/archives.example.edu/
git fetch origin
git pull origin master
