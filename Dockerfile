# Build
FROM golang:alpine AS builder

WORKDIR /build

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 \
    GOOS=linux \
    go build \
    -o /bot \
    ./cmd/main.go

# Deploy
FROM scratch

COPY --from=builder /bot /bin/bot

ENTRYPOINT ["/bin/bot"]