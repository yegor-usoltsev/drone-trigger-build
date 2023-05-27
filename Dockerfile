FROM golang:1.20-alpine AS build
WORKDIR /app
COPY . .
RUN go build

FROM alpine:latest
WORKDIR /app
COPY --from=build /app/drone-trigger-build .
CMD ["/app/drone-trigger-build"]
