FROM golang:1.22 AS builder

WORKDIR /application
# 依賴的包
COPY definition ./definition
COPY tools ./tools
COPY protos ./protos
# 主要的服務
COPY apiserver ./apiserver

RUN cd /application/apiserver && go mod tidy && CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ../service

FROM centos:7
WORKDIR /app
COPY ./apiserver/conf/env.yaml ./conf/env.yaml
COPY --from=builder application/service .
ENTRYPOINT ["./service"]
CMD ["server"]