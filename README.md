Attempt to submit an APM Trace to Datadog without running the `datadog-agent` daemon.

https://docs.datadoghq.com/api/?lang=python#tracing only describes the `tracer` to `datadog-agent` api,
without describing how the `datadog-agent` submits traces to `trace.agent.datadoghq.com`. 

But the [datadog-agent](https://github.com/DataDog/datadog-agent) obviously knows how to submit APM traces.
This application imports the `datadog-agent` go packages as a library,
using the internals for data serialization and submission.

It submits exactly one trace with two hardcoded spans and exits.

All configuration knobs (environment variables and `datadog.yaml`)
are expected to have the same behaviour as in the `datadog-agent` daemon.
The configuration loading logic is called unchanged.

# Warning

This is a very crude initial spike.
It was inspired by https://datadoghq.slack.com/archives/C3SH3KCQG/p1557944796241400.
It may not be a good idea.
It's certainly not a _complete_ or _useful_ implementation.

# Usage

Setup your machine as if the `datadog-agent` would be running on it, including
`api_key` from https://app.datadoghq.com/account/settings

Install go, clone the repo (probably into your `GOPATH`), fetch dependencies.

```shell
$ go version
go version go1.11.5 darwin/amd64
$ dep ensure
```

Run it.

```shell
$ go run cmd/main.go
```

Look for new traces at https://app.datadoghq.com/apm/traces?env=testenv

It compiles down to a 19MB static binary.
That's a bit smaller than the 347MB `datadog-agent` plus the size of a tracing client.

# Testing

There are no tests. This atrocity is 235 lines.
