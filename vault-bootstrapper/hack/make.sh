#!/usr/bin/env bash

goimports -w *.go pkg commands

gofmt -s -w *.go pkg commands