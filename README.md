# Kratos Project Template

## Install Kratos
```
go install github.com/go-kratos/kratos/cmd/kratos/v2@latest
```
## Create a service
```
# Create a template project
kratos new server

cd server
# Add a proto template
kratos proto add api/server/server.proto
# Generate the proto code
kratos proto client api/server/server.proto
# Generate the source code of service by proto file
kratos proto server api/server/server.proto -t internal/service

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

https://go-kratos.dev/en/docs/getting-started/start


docker build -t bunyawat/kratos-k8s-job . 

docker run --rm -p 8000:8000 -p 9000:9000 -v ./configs:/data/conf bunyawat/kratos-k8s-job

docker image ls

docker image push bunyawat/kratos-k8s-job

make build

./bin/kratos-shell-cmd -conf ./configs


failed to create containerd task: 
failed to create shim task: 
OCI runtime create failed: 
runc create failed: 
unable to start container process: 
exec: "kratos-shell-cmd -conf /data/conf": 
stat kratos-shell-cmd -conf /data/conf: 
no such file or directory: 
unknown


helm repo add bitnami https://charts.bitnami.com/bitnami

helm install my-rabbitmq bitnami/rabbitmq --version 12.8.0

kubectl port-forward --namespace default svc/my-rabbitmq 5672:5672

kubectl port-forward --namespace default svc/my-rabbitmq 15672:15672

http://127.0.0.1:15672

echo "Username      : user"

echo "Password      : $(kubectl get secret --namespace default my-rabbitmq -o jsonpath="{.data.rabbitmq-password}" | base64 -d)"

echo "ErLang Cookie : $(kubectl get secret --namespace default my-rabbitmq -o jsonpath="{.data.rabbitmq-erlang-cookie}" | base64 -d)"

