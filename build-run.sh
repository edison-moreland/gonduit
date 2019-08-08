#!/usr/bin/env bash
OUTPUT_DIR="./bin"

if [ -d ${OUTPUT_DIR} ]
then
    rm -rf ${OUTPUT_DIR}
fi

go build -o ${OUTPUT_DIR}/gonduit -v -x