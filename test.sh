#!/bin/bash 

set -ex
set -o pipefail

FOLDERS=$(go list ./... | grep -v /vendor/)
go vet $FOLDERS
echo $FOLDERS | xargs -n1 fgt golint
go test $FOLDERS -v | go-junit-report 