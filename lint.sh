#!/usr/bin/env bash

printf "<<<  gocheckstyle  >>>\n"
gocheckstyle

printf "\n<<<  go-consistent  >>>\n"
go-consistent -v -pedantic ./...

printf "\n<<<  golint strict >>>\n"
golint -min_confidence .3 ./...
