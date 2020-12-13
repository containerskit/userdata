#!/usr/bin/env bash
set -e

if [ -z $1 ]; then
    echo "usage: ${0} EXECUTABLE"
fi

diff -ubBE <(cat tests/help/help.txt) <($1 -help)
