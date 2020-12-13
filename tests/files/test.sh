#!/usr/bin/env bash
set -e

if [ -z $1 ]; then
    echo "usage: ${0} EXECUTABLE"
fi

TMP=tests/files/tmp
mkdir -p $TMP

# trap "rm -rf ${TMP}" EXIT

$1 -path tests/files/userdata.yml
