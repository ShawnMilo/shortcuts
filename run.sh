#!/usr/bin/env bash

go build -o /tmp/shortcuts || exit

/tmp/shortcuts < sample.txt
