FROM golang:1.18-alpine AS builder
WORKDIR /app
COPY . .
RUN go mod download
RUN go build -o ./app ./main.go


FROM alpine:latest AS runner
WORKDIR /app
COPY --from=builder /app .
EXPOSE 8080
ENTRYPOINT ["./app"]