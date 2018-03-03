#!/bin/bash

cd src/cli/main
go build .
mv main ../../../cli
