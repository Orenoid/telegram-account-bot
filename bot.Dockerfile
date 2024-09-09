# Build Stage
FROM golang:1.20-alpine AS build
WORKDIR /app
COPY . .
RUN go build -o telebot ./cmd/telebot/telebot.go
RUN go build -o telebotctl ./cmd/telebotctl/telebotctl.go

# Final Stage
FROM alpine:3.14
RUN apk update && apk add tzdata
WORKDIR /root/
COPY --from=build /app/telebot .
COPY --from=build /app/telebotctl .
COPY bot-entrypoint.sh .
RUN chmod +x bot-entrypoint.sh
ENTRYPOINT ["./bot-entrypoint.sh"]
