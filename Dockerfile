FROM golang:1.24

WORKDIR /app

# Air 설치
RUN go install github.com/air-verse/air@v1.61.7

# 모듈 초기화
COPY go.mod ./
RUN go mod download

# CMD ["./server"]
CMD ["air", "-c", ".air.toml"]