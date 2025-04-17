# 构建阶段
FROM golang:1.23-alpine AS builder

WORKDIR /app

RUN apk add --no-cache git make

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN make wire
RUN make build

# 运行阶段
FROM alpine:latest

RUN apk add --no-cache ca-certificates tzdata

# 设置时区为亚洲/上海
RUN cp /usr/share/zoneinfo/Asia/Shanghai /etc/localtime && echo "Asia/Shanghai" >/etc/timezone


WORKDIR /app

COPY --from=builder /app/sbin/app /app/sbin/app
COPY --from=builder /app/configs /app/configs

# 暴露端口
EXPOSE 8080

# 设置启动命令
RUN cd /app
CMD ["sbin/app", "run"]
