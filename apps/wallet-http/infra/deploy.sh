#!/bin/zsh

cd $(dirname $0)

eval $(minikube docker-env)

docker build -t wallet-http:v1.0.25 ./../src/

kubectl apply -f deployment.yaml
kubectl apply -f service.yaml
kubectl apply -f local-service.yaml
