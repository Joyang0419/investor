# 使用官方 Node.js 鏡像作為基礎
FROM node:22

# 設置工作目錄
WORKDIR /app

# 安裝 migrate-mongo 工具及其依賴的 mongodb 模組
RUN npm install -g migrate-mongo mongodb

# 設置啟動命令，暫不執行任何操作，將在 docker-compose.yml 中指定
CMD ["node"]