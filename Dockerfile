# -- STAGE 1: BUILDER -----------------------------------
FROM golang:1.25.4-alpine AS builder

# Set working directory untuk source code
WORKDIR /app

# Menyalin file dependency
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Menyalin semua source code
COPY . .

# Build aplikasi
# CGO_ENABLED=0 untuk membuat binary statis
# GOOS=linux karena dijalankan di container linux
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main ./cmd

# -- STAGE 2: FINAL IMAGE (SMALL) -----------------------
# Menggunakan image Alpine yang sangat kecil untuk menjalankan binary
FROM alpine:latest

RUN apk --no-cache add ca-certificates

# Set working directory di image final
WORKDIR /root/

# Copy binary yang sudah di-build dari stage 'builder'
COPY --from=builder /app/main .

# User non-root untuk keamanan
RUN adduser -D user1
USER user1

EXPOSE 8086

# Command untuk menjalankan aplikasi
CMD ["./main"]