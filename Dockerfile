# 第一階段：建置階段
FROM golang:1.24.3-alpine AS builder

# 設定工作目錄
WORKDIR /app

# 複製 go.mod 和 go.sum（如果存在）
COPY go.mod go.sum* ./

# 下載依賴
RUN go mod download

# 複製源代碼
COPY . .

# 編譯應用程序
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o war3tool .

# 第二階段：運行階段
FROM alpine:latest

# 安裝必要的運行時依賴（如果需要）
RUN apk --no-cache add ca-certificates

# 設定工作目錄
WORKDIR /app

# 從 builder 階段複製編譯好的二進制文件
COPY --from=builder /app/war3tool .

# 複製所需的資源文件和配置
COPY --from=builder /app/assets ./assets
COPY --from=builder /app/data ./data

# 暴露端口（如果需要修改，請根據實際情況調整）
EXPOSE 8080

# 定義環境變數（運行時需要設定這些）
ENV FA_VALID="YOUR_VALID_KEY"
ENV USERS_FOLDER_PATH="./users"
ENV MAPS_FOLDER_PATH="./maps"
ENV PORT="8080"

# 運行應用程序
CMD ["./war3tool"]
