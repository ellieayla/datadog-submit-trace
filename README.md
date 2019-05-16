Attempt to submit an APM Trace to Datadog without running the `datadog-agent` daemon.

APM Traces are frequently used for transactional web requests.
Spans start, run for a while, have sub-spans, and end with success/failure or some other tags.
That shape is common between an HTTP request (70ms) and a Build (20 minutes).
Maybe APM Traces could be used to model [Lead Time](https://en.wikipedia.org/wiki/Lead_time)
of much longer processes, like a Release.
It could be comprised of spans (all of which have start time and duration)
linked by a common TraceID (the git commit?):

* Git Commit (tag: author, git hash)
* PR Build (tag: build number, build machine)
* PR Test (tag: result)
* PR commentary and merge (tag: people involved)
* Build (tag: build number, build machine)
* Test (tag: result)
* Image Push (tag: image repo/tag)
* Canary Deploy (tag: fqdn?)
* Production Deploy (huzzah)

The measurements for each of those stages would need to be
collected independently from various systems.
This is nearly identical to the concept of
[Distributed Tracing](https://microservices.io/patterns/observability/distributed-tracing.html).
Just on a much longer timescale, and not originating from inside a mostly-static daemon.

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
