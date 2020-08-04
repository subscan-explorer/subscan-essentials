#!/usr/bin/env bash

DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
BASENAME="subscan"
app=${BASENAME}

function build() {
    go build -o ./cmd/subscan -v ./cmd
	echo "Build Success"
}

function install() {
    if [[ "$2" == "" ]]; then
        echo "empty plugin url!"
        return
    fi
    git clone "$1" tmp/"$2"
    if [[ $? -ne 0 ]]; then
        exit $?
    fi
    cd tmp/"$2"
    go build -buildmode=plugin -o ${DIR}/configs/plugins/"$2".so
    if [[ $? -ne 0 ]]; then
        exit $?
    fi
    echo "Build" ${DIR}/configs/plugins "Success"
    rm -rf ${DIR}/tmp
}


function help() {
    echo "usage: ./build.sh build|install[repository][pluginName]"
}

if [[ "$1" == "" ]]; then
    help
elif [[ "$1" == "build" ]];then
    build
elif [[ "$1" == "install" ]];then
    install $2 $3
else
    help
fi
