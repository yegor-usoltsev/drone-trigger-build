# drone-trigger-build

[![Build Status](https://github.com/yegor-usoltsev/drone-trigger-build/actions/workflows/ci.yml/badge.svg)](https://github.com/yegor-usoltsev/drone-trigger-build/actions)
[![Codecov](https://codecov.io/github/yegor-usoltsev/drone-trigger-build/graph/badge.svg?token=W3C1MTB6S4)](https://codecov.io/github/yegor-usoltsev/drone-trigger-build)
[![GitHub Release](https://img.shields.io/github/v/release/yegor-usoltsev/drone-trigger-build?sort=semver)](https://github.com/yegor-usoltsev/drone-trigger-build/releases)
[![Docker Image (docker.io)](https://img.shields.io/docker/v/yusoltsev/drone-trigger-build?label=docker.io&sort=semver)](https://hub.docker.com/r/yusoltsev/drone-trigger-build)
[![Docker Image (ghcr.io)](https://img.shields.io/docker/v/yusoltsev/drone-trigger-build?label=ghcr.io&sort=semver)](https://github.com/yegor-usoltsev/drone-trigger-build/pkgs/container/drone-trigger-build)
[![Docker Image Size](https://img.shields.io/docker/image-size/yusoltsev/drone-trigger-build?sort=semver&arch=amd64)](https://hub.docker.com/r/yusoltsev/drone-trigger-build/tags)

Drone CI / CD plugin to trigger builds for a list of downstream repositories.

This project uses the [Build Create method](https://docs.drone.io/api/builds/build_create/) of the Drone API to
trigger builds with specified parameters for the listed repositories, partially replicating the functionality
of [drone-plugins/drone-downstream](https://github.com/drone-plugins/drone-downstream) but with several improvements:

1. This project creates a new build instead of restarting the previous one.
2. It correctly parses the parameters with commas (e.g. `KEY=VALUE1,VALUE2,VALUE3`).
3. It builds as a multi-platform Docker image (`linux/amd64`, `linux/arm64`).

## Usage

### Drone Pipeline

```yaml
kind: pipeline
name: default

steps:
  - name: trigger
    image: yusoltsev/drone-trigger-build
    settings:
      server: https://drone.example.com
      token:
        from_secret: drone_token
      repositories:
        - octocat/Hello-World
        - octocat/Spoon-Knife
      params:
        - KEY=VALUE
        - FOO=BAR
```

### Docker

```bash
$ docker run --rm \
  -e PLUGIN_SERVER=https://drone.example.com \
  -e PLUGIN_TOKEN=<drone_token> \
  -e PLUGIN_REPOSITORIES=octocat/Hello-World,octocat/Spoon-Knife \
  -e PLUGIN_PARAMS=KEY=VALUE,FOO=BAR \
  -v $(pwd):$(pwd) \
  -w $(pwd) \
  yusoltsev/drone-trigger-build
```

## Docker Images

This application is delivered as a multi-platform Docker image and is available for download from two image registries
of choice: [yusoltsev/drone-trigger-build](https://hub.docker.com/r/yusoltsev/drone-trigger-build)
and [ghcr.io/yegor-usoltsev/drone-trigger-build](https://github.com/yegor-usoltsev/drone-trigger-build/pkgs/container/drone-trigger-build).

## Versioning

This project uses [Semantic Versioning](https://semver.org)

## Contributing

Pull requests are welcome. For major changes,
please [open an issue](https://github.com/yegor-usoltsev/drone-trigger-build/issues/new) first to discuss what you would
like to change. Please make sure to update tests as appropriate.

## License

[MIT](https://github.com/yegor-usoltsev/drone-trigger-build/blob/main/LICENSE)
