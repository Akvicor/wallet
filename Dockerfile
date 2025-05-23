# 构建应用
FROM akvicor/builder:v0.0.9-node20-go22 AS builder
WORKDIR /app
COPY . .

RUN cd frontend && make build
RUN cp -r frontend/build backend/cmd/app/server/web/build
RUN cd backend && make build

# 最小化镜像
FROM debian:12.8-slim
WORKDIR /app

COPY --from=builder /app/backend/bin/wallet ./wallet
COPY --from=builder /app/prod.sh ./prod.sh

RUN ln -sf /usr/share/zoneinfo/Etc/GMT-8 /etc/localtime && \
    mkdir /data && \
    chmod +x ./prod.sh

# healthy check
RUN apt-get update && apt-get install -y curl && rm -rf /var/lib/apt/lists/*

EXPOSE 3000
CMD ["./prod.sh"]

