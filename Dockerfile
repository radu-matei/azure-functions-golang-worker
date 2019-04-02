FROM golang:1.12 as builder
ENV GO111MODULE=on

WORKDIR /go/src/github.com/radu-matei/azure-functions-golang-worker
COPY . .

RUN go build -o workers/go/worker

# compile HTTP Trigger samples that work without any Azure account
WORKDIR /go/src/github.com/radu-matei/azure-functions-golang-worker/sample/HttpTriggerGo
RUN go build -buildmode=plugin -o bin/HttpTriggerGo.so main.go

WORKDIR /go/src/github.com/radu-matei/azure-functions-golang-worker/sample/HttpTriggerHttpResponse
RUN go build -buildmode=plugin -o bin/HttpTriggerHttpResponse.so main.go

# to use a blob-based function you need an Azure storage account and to pass the storage key as env to the container - see readme
# if you have a storage, uncomment the next two steps

#WORKDIR /go/src/github.com/radu-matei/azure-functions-golang-worker/sample/HttpTriggerBlobBindingGo
#RUN go build -buildmode=plugin -o bin/HttpTriggerBlobBindingGo.so main.go

#WORKDIR /go/src/github.com/radu-matei/azure-functions-golang-worker/sample/HttpTriggerBlobBindingInOutGo
#RUN go build -buildmode=plugin -o bin/HttpTriggerBlobBindingInOutGo.so main.go

FROM mcr.microsoft.com/azure-functions/base:2.0

# copy the worker in the pre-defined path
COPY --from=builder /go/src/github.com/radu-matei/azure-functions-golang-worker/workers/go /azure-functions-host/workers/go

# copy all samples
COPY --from=builder /go/src/github.com/radu-matei/azure-functions-golang-worker/sample/ /home/site/wwwroot

ENV AzureWebJobsScriptRoot=/home/site/wwwroot \
    HOME=/home \
    ASPNETCORE_URLS=http://+:80 \
    AZURE_FUNCTIONS_ENVIRONMENT=Development
