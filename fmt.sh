#!/usr/bin/env bash


printf "<<<  Clean module tree  >>>\n"
go mod tidy

printf "<<<  gofumpt  >>>\n"
gofumpt -l -e -s -w .
