#!/bin/sh

. ./help.sh
pacakges="idgen "

for p in $pacakges
do
    build_package $p
done
