perl -pi -e "s/package goberry/package ${1}/g" *.go
perl -pi -e "s/goberry/${1}/" scripts/pythia.sh