# GOTCP L4 scanner
Simple lightweight layer 4 scanner written in golang. The scanner performs two tasks:
* checking connectivity to targets and exposing connectivity status on a server endpoint
* scanning targets over TCP ports range (1- 65,535) and exposing results on a server endpoint
## Pull image from public registry
```sh
docker image pull farazf001/gotcp:latest
```
## Build image locally
```sh
docker image build -t farazf001/gotcp:latest -f Dockerfile .
```
## Run container
```sh
docker container run --name=gotcp -dt -p 8080:8080 farazf001/gotcp:latest
```
## Webserver
The webserver exposes two endpoints with customizable URL parameters:
* ```/health?host=172.18.0.2&port=80```
* ```/report?host=172.18.0.2```
# Sandbox testing
Run the Kind configuration file to provision a Kubernetes cluster:
```sh
kind create cluster --config kind/config.yaml
```
Copy manifest files to Kubernetes control plane:
```sh
docker cp /home/fforoozan/Repositories/gotcp/k8s/ sandbox-control-plane:/gotcp/
```
Deploy Kubernetes API resources:
```sh
docker container exec -it sandbox-control-plane kubectl create ns demo
```
```sh
docker container exec -it sandbox-control-plane kubectl apply -f /gotcp
```
Perform HTTP requests against node ports 30080-30082:
```sh
curl "localhost:30080/health?host=redis.demo.svc&port=6379"
```
```sh
curl "localhost:30080/report?host=redis.demo.svc"
```
Deleting sandbox:
```sh
kind delete cluster --name=sandbox
```