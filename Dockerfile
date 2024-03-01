FROM golang:1.22.0-alpine AS build
WORKDIR /app
COPY . .
RUN go build -ldflags="-s -w"

FROM alpine:latest
WORKDIR /app
COPY --from=build /app/drone-trigger-build .
CMD ["/app/drone-trigger-build"]
