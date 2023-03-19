# Build Stage
FROM golang:1.16-alpine AS build
WORKDIR /app
COPY . .
RUN go build -o telebotctl ./cmd/telebotctl.go

# Final Stage
FROM alpine:3.14
WORKDIR /root/
COPY --from=build /app/telebotctl .
CMD ["./telebotctl"]
