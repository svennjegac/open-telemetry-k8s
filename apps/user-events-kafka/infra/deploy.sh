#!/bin/zsh

cd $(dirname $0)

kubectl apply -f service.yaml
