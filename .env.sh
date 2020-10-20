#!/bin/bash -x

if [[ -z $GOLIBPATH ]]; then
    echo "export GOLIBPATH=path_to_golib"
    exit 1
fi

export GOPATH= # reset $GOPATH

source $GOLIBPATH/env.sh

WD=$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )

export GOPATH="$WD:$GOPATH" # preappend current project into GOPATH
export GOPATH=${GOPATH#:}
export GOPATH=${GOPATH%:}

export PATH=${PATH//:$WD\/bin:/:}
export PATH=${PATH//#$WD\/bin:/}
export PATH=${PATH//%:$WD\/bin/}
export PATH=${PATH}:$WD/bin
[ -f .env ] && source .env || true
