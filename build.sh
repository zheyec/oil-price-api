#!/bin/bash

echo start building

env GOOS=linux GOARCH=amd64 go build -mod=vendor -o lxm-oil-prices main.go
scp lxm-oil-prices root@39.96.21.121:/home/works/chenzheye/lxm-oil-prices
rm lxm-oil-prices