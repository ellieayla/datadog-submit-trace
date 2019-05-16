package submittrace

import (
	"context"
	"fmt"
	upstreamagent "github.com/DataDog/datadog-agent/pkg/trace/agent"
	"github.com/DataDog/datadog-agent/pkg/trace/config"
	"github.com/DataDog/datadog-agent/pkg/trace/flags"
	"github.com/DataDog/datadog-agent/pkg/trace/info"
	"github.com/DataDog/datadog-agent/pkg/trace/osutil"
	"github.com/DataDog/datadog-agent/pkg/trace/pb"
	"github.com/DataDog/datadog-agent/pkg/trace/sampler"

	"github.com/DataDog/datadog-agent/pkg/util/log"
	"os"
	"time"
)

// Run is the entrypoint of our code, which starts the agent.
func Run(ctx context.Context) {
	fmt.Print(info.VersionString())

	cfg, err := config.Load(flags.ConfigPath)
	if err != nil {
		osutil.Exitf("%v", err)
	}
	err = info.InitInfo(cfg) // for expvar & -info option
	if err != nil {
		osutil.Exitf("%v", err)
	}

	if flags.Info {
		if err := info.Info(os.Stdout, cfg); err != nil {
			osutil.Exitf("failed to print info: %s\n", err)
		}
		return
	}

	if err := setupLogger(cfg); err != nil {
		osutil.Exitf("cannot create logger: %v", err)
	}
	defer log.Flush()

	
	// Populate the "agent" hostname
	// cfg.Hostname = "localhost"

	// Lies
	info.Version = "6.5.1"
	info.GitCommit = "1fad9da"

	traceagent := upstreamagent.NewAgent(ctx, cfg)

	log.Trace("Logging at Trace level")

	log.Infof("Enabled? %s", cfg.Enabled)
	log.Infof("Rewritten info %s", info.VersionString())
	log.Infof("X-Datadog-Reported-Languages %s", info.Languages())
	log.Infof("Trace agent running on host %s", cfg.Hostname)
	log.Infof("Endpoints: %s", cfg.Endpoints[0])

	// It's possible that all of these are not needed.
	traceagent.Receiver.Start()
	traceagent.TraceWriter.Start()
	traceagent.StatsWriter.Start()
	traceagent.ServiceMapper.Start()
	traceagent.ServiceWriter.Start()
	traceagent.Concentrator.Start()
	traceagent.ScoreSampler.Start()
	traceagent.ErrorsScoreSampler.Start()
	traceagent.PrioritySampler.Start()
	traceagent.EventProcessor.Start()

	log.Info("Started agent components")

	// Create one trace with two spans
	traceid := time.Now().UnixNano()
	trace := pb.Trace{
		&pb.Span{
			Start:    time.Now().Add(-1 * time.Minute).UnixNano(),
			Duration: time.Duration(5 * time.Second).Nanoseconds(),
			Meta: map[string]string{"datadog.trace_metrics": "true", "env": "testenv",
				"host": cfg.Hostname,
				"lang": "python",
			},
			Metrics: map[string]float64{
				"_sampling_priority_v1": float64(sampler.PriorityUserKeep),
			},
			ParentID: 0,
			TraceID:  uint64(traceid),
			SpanID:   uint64(traceid),
			Service:  "jenkins",
			Name:     "jenkins",
			Resource: "Release",
			Type:     "CICD"},

		&pb.Span{
			TraceID:  uint64(traceid),
			SpanID:   uint64(traceid + 1),
			ParentID: uint64(traceid),
			Service:  "jenkins3",
			Name:     "jenkins3",
			Resource: "Build",
			Meta: map[string]string{
				"build.number": "1234567",
				"env":          "testenv",
			},
			Metrics: map[string]float64{
				"_sampling_priority_v1": float64(sampler.PriorityUserKeep),
			},
			Start:    time.Now().Add(-1 * time.Minute).UnixNano(),
			Duration: time.Duration(3 * time.Second).Nanoseconds(),
		},
	}

	// Transmit
	log.Trace(trace)
	traceagent.Process(trace)

	// Data may have already been sent, or it will be flushed here.
	log.Info("Shutting down agent components, flushing buffered data.")
	traceagent.Concentrator.Stop()
	traceagent.TraceWriter.Stop()
	traceagent.StatsWriter.Stop()
	traceagent.ServiceMapper.Stop()
	traceagent.ServiceWriter.Stop()
	traceagent.ScoreSampler.Stop()
	traceagent.ErrorsScoreSampler.Stop()
	traceagent.PrioritySampler.Stop()
	traceagent.EventProcessor.Stop()
}
