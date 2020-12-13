#!/usr/bin/env bash

APP=$1
if [ -z $APP ]; then
    echo "error: missing required argument APP"
    exit 1
fi

function run() {
    if [ -z $1 ]; then
        echo "usaage: run CASE"
        exit 1
    fi

    output=$(tests/$1/test.sh $APP 2>&1)
    if [ -z "$output" ]; then
        printf "\e[32m%s\tPASSED\e[0m\n" $1
        return 0
    fi

    printf "\e[31m%s\tFAILED\e[0m" $1
    printf "\t\e[31m\n\t"
    echo   "$output"
    printf "\e[0m\n"
    return 1
}

run "help"
run "files"
