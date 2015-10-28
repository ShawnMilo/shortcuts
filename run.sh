#!/usr/bin/env bash

go build || exit

./shortcuts < sample.txt
