FROM golang:1.23.2-alpine AS build
WORKDIR /app
COPY . .
RUN CGO_ENABLED=0 go build -ldflags="-s -w" -trimpath

FROM alpine:latest
WORKDIR /app
COPY --from=build /app/drone-trigger-build .
CMD ["/app/drone-trigger-build"]
