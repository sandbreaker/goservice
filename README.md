# Go API service stub
Go api service stub. Supports basic web/api service framework with ready made utilities such as multi role/app context, static+dynamic configuration, nice logging wrappers, internal metrics with datadog integration, internal alerting system via sentry, basic orm, worker pool, graphql, docker, build wrappers etc.

# Usage
    make run-dev
    Test out on http://localhost:8088/v1/util/version


# Key features
* **Multi app:** Multi app/context support
* **Worker pool:** Basic building blocks for worker pool
* **Service util:** logging, alerts, metrics, dynamic/static config
  

# Progress
* [**done**] project setup
* [**done**] api middleware, logging, build script, etc
* [**done**] static config
* [**done**] internal metric
* [**done**] integration with sentry
* [**done**] basic orm
* [wip] datadog integration (free way)
* [wip] dyanmic config
* [wip] worker pool
* [wip] docker
* [wip] graphql
