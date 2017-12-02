FROM radumatei/functions-runtime:golang

ADD golang-worker /azure-functions-runtime/workers/go/

ENV AzureWebJobsScriptRoot=/home/site/wwwroot
COPY . /home/site/wwwroot 