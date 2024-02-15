FROM golang:1.22 AS builder
COPY . /src
WORKDIR /src
#RUN GOPROXY=https://goproxy.cn make build
RUN make build
COPY "mcafee_certificate.pem" "/tmp/mcafee_certificate.pem"
RUN cat /tmp/mcafee_certificate.pem >> /etc/ssl/certs/ca-certificates.crt

FROM debian:stable-slim
RUN apt-get update && apt-get install -y --no-install-recommends ca-certificates netbase \
    && rm -rf /var/lib/apt/lists/ \
    && apt-get autoremove -y  \
    && apt-get autoclean -y
COPY --from=builder /src/bin /app
WORKDIR /app
EXPOSE 8000
EXPOSE 9000
#VOLUME /data/conf
COPY "configs/config.yaml" "/data/conf/config.yaml"
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
CMD ["./kratos-shell-cmd", "-conf", "/data/conf"]