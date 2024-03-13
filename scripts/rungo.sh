#!/bin/bash

rungo () {
        if [ $# -eq 0 ]
                then npx nodemon --exec go run cmd/main.go 
        elif [ $# -eq 1 ]
                then npx nodemon --exec go run $1 
        fi
}
rungo $1