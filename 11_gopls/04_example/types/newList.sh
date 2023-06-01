#!/bin/sh

TYPE=$1

cat list.tmpl.go | sed -e 's/DUMMYTYPE/'${TYPE}'/g' > ${TYPE}List.go