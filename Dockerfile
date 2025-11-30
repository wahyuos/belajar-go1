# -- STAGE 1: BUILDER -----------------------------------
# Menggunakan image Go yang lebih besar untuk proses kompilasi
FROM golang:1.25-alpine AS builder

# Atur variabel lingkungan untuk static linking (penting untuk image kecil)
ENV CGO_ENABLED=0 GOOS=linux

WORKDIR /app

# Salin go.mod dan go.sum untuk menginstal dependency secara terpisah
COPY go.mod go.sum ./
RUN go mod download

# Salin semua kode sumber
COPY . .

# Kompilasi kode sumber
# -o belajargo1: nama executable yang dihasilkan
# ./cmd: lokasi package main
RUN go build -o /app/belajargo1 ./cmd

# -- STAGE 2: FINAL IMAGE (SMALL) -----------------------
# Menggunakan image Alpine yang sangat kecil untuk menjalankan binary
FROM alpine:latest

# Expose port yang digunakan oleh Gin
EXPOSE 6060

# User non-root untuk keamanan
RUN adduser -D usergo1
USER usergo1

WORKDIR /app

# Salin executable dari stage builder
COPY --from=builder /app/belajargo1 /app/belajargo1

# Jalankan executable
# Anda harus menjalankan container dengan variabel .env dari luar (lihat langkah 3)
CMD ["/app/belajargo1"]