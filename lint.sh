#!/usr/bin/env bash

printf "<<<  gocheckstyle  >>>\n"
gocheckstyle

printf "\n<<<  go-consistent  >>>\n"
go-consistent -v -pedantic ./...
