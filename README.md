# Gotcp L4 scanner
Simple layer 4 scanner written in golang.
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
* ```/healthz?host=172.18.0.2&port=80```
* ```/reportz?host=172.18.0.2&startPort=1&endPort=100000```
