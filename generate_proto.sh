#!/bin/bash
protoc -I protos/ -I ./ --go_out=plugins=grpc:./protos ./protos/*.proto
