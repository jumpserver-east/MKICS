# 第一步：构建 Go 编译环境
FROM golang:1.23.2 AS go-build

ARG GOPROXY=https://goproxy.io
ENV CGO_ENABLED=0 \
    GO111MODULE=on \
    GOPROXY=${GOPROXY} \
    GOPATH=/go \
    PATH=/go/bin:/usr/local/go/bin:$PATH

WORKDIR /opt/evobot

COPY . .

WORKDIR /opt/evobot/cmd/server

RUN go mod download -x \
    && go build -o /opt/evobot/evobot ./main.go


# 第二步：构建前端
FROM node:22.11.0 AS frontend-build

ARG NPM_REGISTRY="https://registry.npmmirror.com"
ENV NPM_REGISTRY=${NPM_REGISTRY}

WORKDIR /app

COPY frontend ./frontend

RUN npm config set registry ${NPM_REGISTRY} \
    && yarn config set registry ${NPM_REGISTRY} \
    && cd frontend \
    && yarn install \
    && yarn build


# 第三步：精简最终产物
FROM alpine:latest AS final

RUN apk add --no-cache bash

WORKDIR /opt/evobot

# 拷贝可执行文件与配置
COPY --from=go-build /opt/evobot/evobot .
COPY --from=go-build //opt/evobot/evobot/cmd/server/conf/config-example.yaml ./config.yaml

# 拷贝前端静态资源（如有嵌入）
COPY --from=frontend-build /app/frontend/dist ./frontend/dist

EXPOSE 24916

CMD ["/opt/evobot/evobot", "-c", "/opt/evobot/config.yaml"]
