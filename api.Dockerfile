# Build Stage
FROM golang:1.20-alpine AS build
WORKDIR /app
COPY . .
RUN go build -o api ./cmd/api/api.go

# Final Stage
FROM alpine:3.14
RUN apk update && apk add tzdata
WORKDIR /root/
COPY --from=build /app/api .
COPY api-entrypoint.sh .
RUN chmod +x api-entrypoint.sh
ENTRYPOINT ["./api-entrypoint.sh"]
