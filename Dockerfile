# Build Stage
FROM golang:1.17-alpine AS build
WORKDIR /app
COPY . .
RUN go build -o telebotctl ./cmd/telebotctl.go
RUN go build -o migrate-cli ./cmd/migrate_cli.go

# Final Stage
FROM alpine:3.14
WORKDIR /root/
COPY --from=build /app/telebotctl .
COPY entrypoint.sh .
RUN chmod +x entrypoint.sh
ENTRYPOINT ["./entrypoint.sh"]
