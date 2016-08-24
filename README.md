# ckeeper

[![Build Status](https://travis-ci.org/yeasy/ckeeper.svg?branch=master)](https://travis-ci.org/yeasy/ckeeper)
[![Go Report Card](https://goreportcard.com/badge/github.com/yeasy/ckeeper)](https://goreportcard.com/report/github.com/yeasy/ckeeper)
Keep the health of container applications.

ckeeper can be deployed on each docker host, and it will automatically check the health status, e.g., validate if the web service online, and run specified operations on `unhealthy` containers, e.g., `restart`.

The default configuration is in [ckeeper.yaml](ckeeper.yaml).

## Features

* Automatically watch the health status of application as a nanny.
* Support container selection based on flexible filtering options.
* Support customized healing actions for unhealthy applications.


## Usage

There are two ways to use ckeeper: docker image (recommended) and local installation.

### Run in container
```sh
$ docker run --rm \
	 --name ckeeper \
	 yeasy/ckeeper \
	 start --logging-level=debug
```

### Local run

```sh
$ make run
```

## Configuration File

### logging

* `level`: Set output logging message level: `debug|info|warning|error|critical`.

### host

* `daemon`: The daemon url of the container host

### check

* `retries`: Each check will run how many tries before thinking it unhealthy.
* `interval`: Seconds between two check.


### Rules

Define how to decide unhealthy status and what action to take.

You can put multiple rules here, each rule will be processed in order.

For each rule:

* `option`: Used to select containers to check health, following the [ListOption](https://docs.docker.com/engine/reference/api/docker_remote_api_v1.24/#list-containers).
* `target`: Think as healthy if the target cmd's running result is true.
* `action`: What action to take if the target is not met. If the target is not set, then take action by default.

The target cmd can be customized easily to achieve flexible health as the user needs.

E.g., the following rule will keep every running web containers serving normally on port 80, otherwise, restart it.

```yaml
  rule_web:
    option:
      All: true
      Filters:
        status:
          - running
    target: '[[ `curl -sL -w "%{http_code}\\n" CONTAINER:80 -o /dev/null` == "200" ]]'
    action: "restart"
```

The following rule will start all containers in `exited` or `paused` status.

```yaml
  rule_start:
    option:   # container filter
      All: true
      Filters:
        status:
          - exited
          - paused
    action: "start"
```

## Testing

```sh
$ make test
```

## TODO

* Add more docs and test cases.
* Design more powerful rule engine.
