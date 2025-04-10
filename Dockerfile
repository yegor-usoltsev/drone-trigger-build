FROM alpine:latest
ENTRYPOINT ["/drone-trigger-build"]
COPY drone-trigger-build /
