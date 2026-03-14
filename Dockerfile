FROM alpine:latest
RUN apk add --no-cache --update ca-certificates tzdata
ARG TARGETPLATFORM
ENTRYPOINT ["/drone-trigger-build"]
COPY $TARGETPLATFORM/drone-trigger-build /
