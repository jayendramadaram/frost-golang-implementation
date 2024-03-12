#!/bin/bash

rungo () {
        if [ $# -eq 0 ]
                then npx nodemon --exec go run cmd/main.go --signal SIGTERM
        elif [ $# -eq 1 ]
                then npx nodemon --exec go run $1 --signal SIGTERM
        fi
}
rungo $1