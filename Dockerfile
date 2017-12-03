FROM golang as builder

WORKDIR /go/src/github.com/radu-matei/azure-functions-golang-worker
COPY . .

RUN go build -o golang-worker

WORKDIR /go/src/github.com/radu-matei/azure-functions-golang-worker/sample/HttpTriggerGo
RUN go build -buildmode=plugin -o bin/HttpTriggerGo.so main.go


FROM radumatei/functions-runtime:golang

COPY --from=builder /go/src/github.com/radu-matei/azure-functions-golang-worker/golang-worker /azure-functions-runtime/workers/go/
COPY --from=builder /go/src/github.com/radu-matei/azure-functions-golang-worker/sample/ /home/site/wwwroot

ENV AzureWebJobsScriptRoot=/home/site/wwwroot
