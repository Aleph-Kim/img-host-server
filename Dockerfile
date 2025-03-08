FROM golang:1.20-alpine

WORKDIR /app
COPY . .

# 모듈 초기화
RUN go mod init img-host || true
RUN go mod tidy

RUN go build -o server .

EXPOSE 3000

CMD ["./server"]
