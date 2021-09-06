# GOTCP L4 scanner
Simple lightweight layer 4 scanner written in golang. The scanner performs two tasks:
* checking connectivity to targets and exposing connectivity status on a server endpoint
* scanning targets over TCP ports range (1- 65,535) and exposing results on a server endpoint
# Image pull
```sh
docker image pull farazf001/gotcp:latest
```
## Image build
```sh
docker build -t farazf001/gotcp:latest -f Dockerfile .
```
## Container run
```sh
docker container run --name=gotcp -dt -p 8080:8080 farazf001/gotcp:latest
```
## Webserver
The webserver exposes two endpoints with customizable URL parameters:
* ```/health?host=172.18.0.2&port=80```
* ```/report?host=172.18.0.2```
## Sandbox Kubernetes cluster
Run the Kind configuration file to provision a cluster:
```sh
kind create cluster --config kind/config.yaml
```
Copy project files to and from Kubernetes control plane:
```sh
docker cp /home/fforoozan/repos/gotcp/k8s/ sandbox-control-plane:/gotcp/
```
```sh
docker cp sandbox-control-plane:/gotcp/k8s/ /home/fforoozan/repos/gotcp/
```