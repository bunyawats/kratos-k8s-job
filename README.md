# Kratos Project Template

## Install Kratos
```
go install github.com/go-kratos/kratos/cmd/kratos/v2@latest
```
## Create a service
```
# Create a template project
kratos new kratos-k8s-job

cd server
# Add a proto template
% kratos proto add api/scheduler/v1/job.proto  
# Generate the proto code
kratos proto client api/scheduler/v1/job.proto  
# Generate the source code of service by proto file
kratos proto server api/scheduler/v1/job.proto -t internal/service

go generate ./...
go build -o ./bin/ ./...
./bin/server -conf ./configs
```
## Generate other auxiliary files by Makefile
```
# Download and update dependencies
make init
# Generate API files (include: pb.go, http, grpc, validate, swagger) by proto file
make api
# Generate all files
make all
```
## Automated Initialization (wire)
```
# install wire
go get github.com/google/wire/cmd/wire

# generate wire
cd cmd/server
wire
```

## Docker
```bash
# build
docker build -t <your-docker-image-name> .

# run
docker run --rm -p 8000:8000 -p 9000:9000 -v </path/to/your/configs>:/data/conf <your-docker-image-name>
```

===============================================================

https://go-kratos.dev/en/docs/getting-started/start
* make build
* ./bin/kratos-shell-cmd -conf ./configs


## Update Config data model
* make config

## SQLC
https://sqlc.dev/
* cd internal/data/mysql
* sqlc generate

## Docker
* docker build -t bunyawat/kratos-k8s-job .
* docker run --rm -p 8000:8000 -p 9000:9000  bunyawat/kratos-k8s-job
* docker image ls
* docker image push bunyawat/kratos-k8s-job


## RabbitMQ
https://rabbitmq.com/

* helm repo add bitnami https://charts.bitnami.com/bitnami
* helm install my-rabbitmq bitnami/rabbitmq --version 12.8.0
* kubectl port-forward --namespace default svc/my-rabbitmq 5672:5672
* kubectl port-forward --namespace default svc/my-rabbitmq 15672:15672
rabbitmq dashboard: http://127.0.0.1:15672

* echo "Username      : user"  
* echo "Password      : $(kubectl get secret --namespace default my-rabbitmq -o jsonpath="{.data.rabbitmq-password}" | base64 -d)"  
* echo "ErLang Cookie : $(kubectl get secret --namespace default my-rabbitmq -o jsonpath="{.data.rabbitmq-erlang-cookie}" | base64 -d)"  


## Kubernetes
* kind get clusters  
* export POD_NAME=$(kubectl get pods -n kubernetes-dashboard -l "app.kubernetes.io/name=kubernetes-dashboard,app.kubernetes.io/instance=kubernetes-dashboard" -o jsonpath="{.items[0].metadata.name}")  
* kubectl -n kubernetes-dashboard port-forward $POD_NAME 8443:8443
Kubernetes dashboard: https://127.0.0.1:8443/
* kubectl create -f cronjob.yaml
* kubectl get cronjob kratos-k8s-job
* kubectl delete cronjob kratos-k8s-job


## InfluxDB
```azure
SELECT "/memory/classes/metadata/other:bytes", "time"
FROM "metrics"
WHERE
time >= now() - interval '12 hours'
```

## READ Configuration for system environment 
export KRATOS_INFLUXDB_TOKEN=<your secret token>
