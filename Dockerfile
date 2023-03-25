# Build Stage
FROM golang:1.20-alpine AS build
WORKDIR /app
COPY . .
RUN go build -o telebot ./cmd/telebot/telebot.go
RUN go build -o telebotctl ./cmd/telebotctl/telebotctl.go

# Final Stage
FROM alpine:3.14
WORKDIR /root/
COPY --from=build /app/telebot .
COPY --from=build /app/telebotctl .
COPY entrypoint.sh .
RUN chmod +x entrypoint.sh
ENTRYPOINT ["./entrypoint.sh"]
