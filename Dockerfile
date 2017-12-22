#start from golang 1.10-rc (as multiple plugins with same package name fails in golang-1.9.x)
FROM golang:1.10-rc as builder

WORKDIR /go/src/github.com/radu-matei/azure-functions-golang-worker
COPY . .


RUN go get -u github.com/golang/dep/...
RUN dep ensure

RUN go build -o golang-worker

# compile HTTP Trigger sample that works without any Azure account
WORKDIR /go/src/github.com/radu-matei/azure-functions-golang-worker/sample/HttpTriggerGo
RUN go build -buildmode=plugin -o bin/HttpTriggerGo.so main.go


# to use a blob-based function you need an Azure storage account and to pass the storage key as env to the container - see readme
# if you have a storage, uncomment the next two steps

#WORKDIR /go/src/github.com/radu-matei/azure-functions-golang-worker/sample/HttpTriggerBlobBindingGo
#RUN go build -buildmode=plugin -o bin/HttpTriggerBlobBindingGo.so main.go

#WORKDIR /go/src/github.com/radu-matei/azure-functions-golang-worker/sample/HttpTriggerBlobBindingInOutGo
#RUN go build -buildmode=plugin -o bin/HttpTriggerBlobBindingInOutGo.so main.go

# this is just the Azure Functions Runtime configured to recognize .go functions and to start the worker
FROM radumatei/functions-runtime:golang

# copy the worker in the pre-defined path
COPY --from=builder /go/src/github.com/radu-matei/azure-functions-golang-worker/golang-worker /azure-functions-runtime/workers/go/

# copy all samples
COPY --from=builder /go/src/github.com/radu-matei/azure-functions-golang-worker/sample/ /home/site/wwwroot

ENV AzureWebJobsScriptRoot=/home/site/wwwroot
ENV ASPNETCORE_URLS=http://+:80
