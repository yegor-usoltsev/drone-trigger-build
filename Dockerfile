FROM alpine:latest
RUN apk add --no-cache --update ca-certificates tzdata
ENTRYPOINT ["/drone-trigger-build"]
COPY drone-trigger-build /
