# --- 1. AŞAMA: Derleme (Build) ---
FROM golang:1.25-alpine AS builder

# Çalışma dizinini ayarla
WORKDIR /app

# Bağımlılıkları kopyala ve yükle
COPY go.mod  go.sum ./
RUN go mod download

# Tüm kodları kopyala
COPY . .

# Uygulamayı derle (CGO_ENABLED=0 statik bir binary üretir)
RUN CGO_ENABLED=0 GOOS=linux go build -o main .

# --- 2. AŞAMA: Çalıştırma (Final) ---
FROM alpine:latest

WORKDIR /root/

# Sadece derlenmiş binary dosyasını builder aşamasından al
COPY --from=builder /app/main .

# Uygulamanın çalışacağı portu belirt
EXPOSE 8080

# Uygulamayı başlat
CMD ["./main"]