Azure Functions Golang Worker
=============================

This project aims to add Golang support for Azure Functions.


How to run the sample
---------------------

In order to register Golang as a worker for the Azure Functions Runtime you need to [implement an `IWorkerProvider` as described here](https://github.com/Azure/azure-webjobs-sdk-script/wiki/Language-Extensibility).
I already did this in a [fork of the Azure Functions Runtime and you can find all modifications here](https://github.com/Azure/azure-webjobs-sdk-script/compare/dev...radu-matei:golang-worker) and pushed a Docker image on Docker Hub based on the [Dockerfile here](https://github.com/radu-matei/azure-webjobs-sdk-script/blob/golang-worker/Dockerfile)

To build the worker you need to: 

- `make golang-worker` , then go to the [samples folder](/sample) and:
 
- `docker build -t azure-functions-go-sample .`
 
- `docker run -e "ASPNETCORE_URLS=http://+:80" -p 81:80 -it azure-functions-go-sample`



Disclaimer
----------
This is not an official Azure Project - it is an unofficial project to support native Golang in Azure Functions by implementing the Worker for v2 - [more details here](https://github.com/Azure/azure-webjobs-sdk-script/wiki/Language-Extensibility)

It is not officially supported by Microsoft and it is not guaranteed to be supported or even work.