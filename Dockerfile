FROM alpine:latest
RUN apk add --no-cache --update ca-certificates tini tzdata
ARG TARGETPLATFORM
ENTRYPOINT ["tini", "--"]
CMD ["/drone-trigger-build"]
COPY $TARGETPLATFORM/drone-trigger-build /
