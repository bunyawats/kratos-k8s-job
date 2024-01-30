FROM golang:1.21 AS builder

COPY . /src
WORKDIR /src

#RUN GOPROXY=https://goproxy.cn make build
RUN make build

FROM debian:stable-slim

RUN apt-get update && apt-get install -y --no-install-recommends \
		ca-certificates  \
        netbase \
        && rm -rf /var/lib/apt/lists/ \
        && apt-get autoremove -y && apt-get autoclean -y

COPY --from=builder /src/bin /app

WORKDIR /app

EXPOSE 8000
EXPOSE 9000
#VOLUME /data/conf

COPY "configs/config.yaml" "/data/conf/config.yaml"

CMD ["./kratos-shell-cmd", "-conf", "/data/conf"]