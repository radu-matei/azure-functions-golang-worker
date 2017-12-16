FROM golang:1.10-rc as builder
#FROM golang as builder

WORKDIR /go/src/github.com/radu-matei/azure-functions-golang-worker
COPY . .


#RUN go get -u github.com/golang/dep/...
#RUN dep ensure

RUN go build -o golang-worker

WORKDIR /go/src/github.com/radu-matei/azure-functions-golang-worker/sample/HttpTriggerGo
RUN go build -buildmode=plugin -o bin/HttpTriggerGo.so main.go

#WORKDIR /go/src/github.com/radu-matei/azure-functions-golang-worker/sample/BlobTriggerGo
#RUN go build -buildmode=plugin -o bin/BlobTriggerGo.so main.go

FROM radumatei/functions-runtime:golang

COPY --from=builder /go/src/github.com/radu-matei/azure-functions-golang-worker/golang-worker /azure-functions-runtime/workers/go/
COPY --from=builder /go/src/github.com/radu-matei/azure-functions-golang-worker/sample/ /home/site/wwwroot

ENV AzureWebJobsScriptRoot=/home/site/wwwroot
ENV ASPNETCORE_URLS=http://+:80
