FROM golang:1.21.5-alpine AS build
WORKDIR /app
COPY . .
RUN go build -ldflags="-s -w"

FROM alpine:latest
WORKDIR /app
COPY --from=build /app/drone-trigger-build .
CMD ["/app/drone-trigger-build"]
