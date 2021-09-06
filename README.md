# GOTCP L4 scanner
Simple lightweight layer 4 scanner written in golang. The scanner performs two tasks:
* checking connectivity to targets and exposing connectivity status on a server endpoint
* scanning targets over TCP ports range (1- 65,535) and exposing results on a server endpoint
## Pull image from Dockerhub
```sh
docker image pull farazf001/gotcp:latest
```
## Build image locally
```sh
docker build -t farazf001/gotcp:latest -f Dockerfile .
```
## Run container
```sh
docker container run --name=gotcp -dt -p 8080:8080 farazf001/gotcp:latest
```
## Webserver endpoints
The webserver exposes two endpoints with customizable URL parameters:
* ```/health?host=172.18.0.2&port=80```
* ```/report?host=172.18.0.2```
